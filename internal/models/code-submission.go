package models

import (
	"context"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionSubmission struct {
	ID         string `json:"id" bson:"_id"`
	QuestionID string `json:"question_id" bson:"question_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	Status     string `json:"status" bson:"status"` // "success" or "fail"
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func CreateSubmission(questionSubmission *QuestionSubmission) (*mongo.InsertOneResult, error) {
	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CODE_SUBMISSION_COLLECTION).InsertOne(context.TODO(), bson.M{
		"question_id": questionSubmission.QuestionID,
		"user_id": questionSubmission.UserID,
		"status": questionSubmission.Status,
		"created_at": questionSubmission.CreatedAt,
	})
	return result, err
}