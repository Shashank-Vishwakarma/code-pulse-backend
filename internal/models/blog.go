package models

import (
	"context"
	"strings"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Blog struct {
	ID              string    `json:"id" bson:"_id"`
	Title           string    `json:"title" bson:"title"`
	Body            string    `json:"body" bson:"body"`
	ImageURL        string    `json:"imageUrl" bson:"imageUrl"`
	IsBlogPublished bool      `json:"isBlogPublished" bson:"isBlogPublished"`
	Slug            string    `json:"slug" bson:"slug"`
	Comments        []Comment `json:"comments,omitempty" bson:"comments"`
	UpVotes         []string  `json:"upVotes,omitempty" bson:"upVotes"`
	DownVotes       []string  `json:"downVotes,omitempty" bson:"downVotes"`
	AuthorID        primitive.ObjectID    `json:"authorId" bson:"authorId"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
}

func CreateBlog(blog *Blog) (*mongo.InsertOneResult, error) {
	words := strings.Split(blog.Title, " ")
	slug := strings.Join(words, "-")

	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.BLOG_COLLECTION).InsertOne(context.TODO(), bson.M{
		"title":           blog.Title,
		"body":            blog.Body,
		"imageUrl":        blog.ImageURL,
		"isBlogPublished": blog.IsBlogPublished,
		"slug":            slug,
		"authorId":        blog.AuthorID,
		"createdAt":       time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
