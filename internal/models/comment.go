package models

import "time"

type Comment struct {
	ID        string    `json:"id" bson:"_id"`
	Body      string    `json:"body" bson:"body"`
	UserID    string    `json:"userId" bson:"userId"`
	BlogID    string    `json:"blogId" bson:"blogId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
