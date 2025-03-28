package models

type ChallengeData struct {
	Question string `json:"question" bson:"question"`
	Answer   string `json:"answer" bson:"answer"`
}

type Challenge struct {
	ID        string          `json:"id" bson:"_id"`
	Title     string          `json:"title" bson:"title"`
	Data      []ChallengeData `json:"data" bson:"data"`
	CreatedAt string          `json:"created_at" bson:"created_at"`
}