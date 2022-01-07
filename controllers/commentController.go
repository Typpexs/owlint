package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"owlint/models"
	"owlint/utilities"
	"strconv"
	"time"

	"owlint/config"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

func WriteSlack(ctx context.Context, newComment models.NewComment) {
	commentCollection := config.MI.DB.Collection("comment")

	apiSlack := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	_, timestamp, err := apiSlack.PostMessage(
		os.Getenv("SLACK_URL_CHANNEL"),
		slack.MsgOptionText(newComment.TextFr, false),
	)

	if err != nil {
		fmt.Println(err)
	}

	commentCollection.FindOneAndUpdate(ctx, bson.M{"targetId": newComment.TargetId},
		bson.M{
			"$set": bson.M{"tsSlack": timestamp},
		})

}

func UpdateCommentReplies(ctx context.Context, newComment models.NewComment) *mongo.UpdateResult {
	commentCollection := config.MI.DB.Collection("comment")
	oldComment := new(models.Comment)
	err := commentCollection.FindOne(ctx, bson.M{"targetId": newComment.TargetId}).Decode(&oldComment)
	if err != nil {
		return nil
	}

	update := bson.M{
		"$push": bson.M{"replies": newComment},
	}

	result, err := commentCollection.UpdateOne(ctx, bson.M{"targetId": newComment.TargetId}, update)
	if err != nil {
		return nil
	}

	apiSlack := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	_, _, err = apiSlack.PostMessage(
		os.Getenv("SLACK_URL_CHANNEL"),
		slack.MsgOptionText(newComment.TextFr, false),
		slack.MsgOptionTS(oldComment.TsSlack),
	)

	if err != nil {
		fmt.Println("error : ", err)
	}

	return result
}

func AddNewComment(c *fiber.Ctx) error {
	commentCollection := config.MI.DB.Collection("comment")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	newComment := new(models.NewComment)

	if err := c.BodyParser(newComment); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	newComment.PublishedAt = strconv.FormatInt(time.Now().Unix(), 10)
	if newComment.TextEn == "" {
		translated, err := utilities.TranslateText("fr", "en", newComment.TextFr)
		if err != nil {
			fmt.Println(err)
		}
		newComment.TextEn = translated
	}

	if newComment.TextFr == "" {
		translated, err := utilities.TranslateText("en", "fr", newComment.TextEn)
		if err != nil {
			fmt.Println(err)
		}
		newComment.TextFr = translated
	}

	_, err := commentCollection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "targetId", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		fmt.Println(err)
	}

	newComment.TargetId = c.Params("id")

	result, err := commentCollection.InsertOne(ctx, newComment)

	if err != nil {
		if IsDup(err) {
			result := UpdateCommentReplies(ctx, *newComment)
			if result == nil {
				return c.Status(500).JSON(fiber.Map{
					"success": false,
					"message": "NewComment failed to insert to parent",
					"error":   err,
				})
			}
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"description": "Comment created",
				"content":     newComment,
			})
		} else {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "NewComment failed to insert",
				"error":   err,
			})
		}
	}

	WriteSlack(ctx, *newComment)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"description": "Comment created",
		"content":     result,
	})

}

func GetComment(c *fiber.Ctx) error {
	commentCollection := config.MI.DB.Collection("comment")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var comments models.Comment
	err := commentCollection.FindOne(ctx, bson.M{"targetId": c.Params("id")}).Decode(&comments)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"description": "comment not found",
		})
	}

	var allComments []models.Comment

	if comments.Replies != nil {
		for index := 0; index < len(comments.Replies); index++ {
			allComments = append(allComments, comments.Replies[index])
		}
		comments.Replies = nil
		allComments = append(allComments, comments)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"item":        allComments,
		"description": "comments matching targetId",
	})
}
