package models

import (
	"context"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID                        string    `json:"id" bson:"_id"`
	Name                      string    `json:"name" bson:"name"`
	Username                  string    `json:"username" bson:"username"`
	Email                     string    `json:"email" bson:"email"`
	Password                  string    `json:"password,omitempty" bson:"password"`
	IsEmailVerified           bool      `json:"is_verified" bson:"is_verified"`
	VerificationCode          string    `json:"verification_code" bson:"verification_code"`
	VerificationCodeExpiresAt time.Time `json:"verification_code_expires_at" bson:"verification_code_expires_at"`
	CreatedAt                 time.Time `json:"created_at" bson:"created_at"`
}

func CreateUser(user *User) (*mongo.InsertOneResult, error) {
	result, err := database.UserCollection.InsertOne(context.Background(), bson.M{
		"name":                         user.Name,
		"username":                     user.Username,
		"email":                        user.Email,
		"password":                     user.Password,
		"is_email_verified":            false,
		"verification_code":            user.VerificationCode,
		"verification_code_expires_at": time.Now().Add(time.Hour),
		"created_at":                   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
