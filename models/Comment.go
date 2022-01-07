package models

type Comment struct {
	TextFr      string    `json:"textFr" bson:"textFr"`
	TextEn      string    `json:"textEn" bson:"textEn"`
	PublishedAt string    `json:"publishedAt" bson:"publishedAt"`
	AuthorId    string    `json:"authorId" bson:"authorId"`
	TargetId    string    `json:"targetId" bson:"targetId"`
	Replies     []Comment `json:"replies" bson:"replies"`
	TsSlack     string    `json:"tsSlack" bson:"tsSlack"`
}

type NewComment struct {
	TextFr      string `json:"textFr" bson:"textFr"`
	TextEn      string `json:"textEn" bson:"textEn"`
	PublishedAt string `json:"publishedAt" bson:"publishedAt"`
	AuthorId    string `json:"authorId" bson:"authorId"`
	TargetId    string `json:"targetId" bson:"targetId"`
}
