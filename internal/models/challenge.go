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

type ChallegeQuestion struct {
	Question string `json:"question"`
	Options []string `json:"options"`
	CorrectAnswer string `json:"correct_answer"`
}

type Challenge struct {
	ID        string          `json:"id" bson:"_id"`
	Topic     string          `json:"topic" bson:"topic"`
	Data      []ChallegeQuestion `json:"data" bson:"data"`
	UserID    string          `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func CreateChallenge(challenge *Challenge) (*mongo.InsertOneResult, error) {
	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).InsertOne(context.TODO(), &bson.M{
		"topic": challenge.Topic,
		"data": challenge.Data,
		"user_id": challenge.UserID,
		"created_at": time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}