package models

import (
	"time"
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
	Cpp        Language = "C++"
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
	ID                  string         `json:"id" bson:"_id,omitempty"`
	Title               string         `json:"title" bson:"title" validate:"required"`
	Slug                string         `json:"slug" bson:"slug" validate:"required"`
	Description         string         `json:"description" bson:"description" validate:"required"`
	Difficulty          Difficulty     `json:"difficulty" bson:"difficulty" validate:"required"`
	Tags                []string       `json:"tags" bson:"tags"`
	Companies           []string       `json:"companies,omitempty" bson:"companies,omitempty"`
	Hints               []string       `json:"hints,omitempty" bson:"hints,omitempty"`
	TestCases           []TestCase     `json:"testCases" bson:"testCases"`
	CodeSnippets        []CodeSnippet  `json:"codeSnippets,omitempty" bson:"codeSnippets,omitempty"`
	TotalSubmissions    int            `json:"totalSubmissions,omitempty" bson:"totalSubmissions,omitempty"`
	Status              QuestionStatus `json:"status" bson:"status"`
	IsQuestionPublished bool           `json:"isQuestionPublished" bson:"isQuestionPublished"`
	Likes               []string       `json:"likes,omitempty" bson:"likes,omitempty"`
	Dislikes            []string       `json:"dislikes,omitempty" bson:"dislikes,omitempty"`
	AuthorID            string         `json:"authorId,omitempty" bson:"authorId,omitempty"`
	CreatedAt           time.Time      `json:"createdAt" bson:"createdAt"`
}
