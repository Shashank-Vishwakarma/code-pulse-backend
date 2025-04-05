package models

import (
	"context"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	ID        primitive.ObjectID    `json:"id" bson:"_id"`
	Body      string    `json:"body" bson:"body"`
	UserID    primitive.ObjectID    `json:"userId" bson:"userId"`
	BlogID    primitive.ObjectID     `json:"blogId" bson:"blogId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func InsertDocumentInComments(comment *Comment) (*mongo.InsertOneResult, error) {
	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.COMMENT_COLLECTION).InsertOne(
		context.TODO(),
		bson.M{
			"body": comment.Body,
			"userId": comment.UserID,
			"blogId": comment.BlogID,
			"createdAt": time.Now(),
		},
	)

	return result, err
}