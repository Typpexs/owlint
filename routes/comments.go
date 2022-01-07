package routes

import (
	"owlint/controllers"

	"github.com/gofiber/fiber/v2"
)

func CommentsRoute(route fiber.Router) {
	route.Get("/:id/comments", controllers.GetComment)
	route.Post("/:id/comments", controllers.AddNewComment)
}
