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

type UserSelectedAnswer struct {
	Question string `json:"question" bson:"question"`
	Answer string `json:"answer" bson:"answer"`
}

type ChallegeQuestion struct {
	Question string `json:"question" bson:"question"`
	Options []string `json:"options" bson:"options"`
	CorrectAnswer string `json:"correct_answer" bson:"correct_answer"`
}

type Challenge struct {
	ID        string          `json:"id" bson:"_id"`
	Title string 	`json:"title" bson:"title"`
	Topic     string          `json:"topic" bson:"topic"`
	Difficulty string		  `json:"difficulty" bson:"difficulty"`
	Data      []ChallegeQuestion `json:"data" bson:"data"`
	Score string `json:"score" bson:"score"`
	UserSelectedAnswers []UserSelectedAnswer `json:"user_selected_answers" bson:"user_selected_answers"`
	UserID    primitive.ObjectID          `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func CreateChallenge(challenge *Challenge) (*mongo.InsertOneResult, error) {
	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).InsertOne(context.TODO(), &bson.M{
		"title": challenge.Title,
		"topic": challenge.Topic,
		"difficulty": challenge.Difficulty,
		"data": challenge.Data,
		"user_id": challenge.UserID,
		"created_at": time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}