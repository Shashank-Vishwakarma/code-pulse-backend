package models

import (
	"context"
	"strings"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Difficulty string

type QuestionStatus string

type Language string

const (
	Easy   Difficulty = "Easy"
	Medium Difficulty = "Medium"
	Hard   Difficulty = "Hard"

	Pending  QuestionStatus = "Pending"
	Approved QuestionStatus = "Approved"
	Rejected QuestionStatus = "Rejected"

	Python     Language = "Python"
	JavaScript Language = "JavaScript"
)

type TestCase struct {
	Input       string `json:"input" bson:"input"`
	Output      string `json:"output" bson:"output"`
	Explanation string `json:"explanation,omitempty" bson:"explanation,omitempty"`
}

type CodeSnippet struct {
	Language Language `json:"language" bson:"language"`
	Code     string   `json:"code" bson:"code"`
}

type Question struct {
	ID           string         `json:"id" bson:"_id,omitempty"`
	Title        string         `json:"title" bson:"title"`
	Slug         string         `json:"slug" bson:"slug"`
	Description  string         `json:"description" bson:"description"`
	Difficulty   Difficulty     `json:"difficulty" bson:"difficulty"`
	Tags         []string       `json:"tags" bson:"tags"`
	Companies    []string       `json:"companies,omitempty" bson:"companies,omitempty"`
	Hints        []string       `json:"hints,omitempty" bson:"hints,omitempty"`
	TestCases    []TestCase     `json:"testCases" bson:"testCases"`
	CodeSnippets []CodeSnippet  `json:"codeSnippets,omitempty" bson:"codeSnippets,omitempty"`
	Status       QuestionStatus `json:"status" bson:"status"`
	AuthorID     string         `json:"authorId,omitempty" bson:"authorId,omitempty"`
	CreatedAt    time.Time      `json:"createdAt" bson:"createdAt"`
}

func CreateQuestion(q *Question) (*mongo.InsertOneResult, error) {
	words := strings.Split(q.Title, " ")
	slug := strings.Join(words, "-")

	result, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).InsertOne(context.TODO(), bson.M{
		"title":        q.Title,
		"slug":         slug,
		"description":  q.Description,
		"difficulty":   q.Difficulty,
		"tags":         q.Tags,
		"companies":    q.Companies,
		"hints":        q.Hints,
		"testCases":    q.TestCases,
		"codeSnippets": q.CodeSnippets,
		"status":       Pending,
		"authorId":     q.AuthorID,
		"createdAt":    time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
