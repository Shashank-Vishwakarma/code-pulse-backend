package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/queue"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context) {
	var body request.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// validate the request body
	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error validating the request body", nil)
		return
	}

	if body.Password != body.ConfirmPassword {
		logrus.Error("Passwords do not match")
		response.HandleResponse(c, http.StatusBadRequest, "Passwords do not match", nil)
		return
	}

	usernameExistsWithUsername := database.UserCollection.FindOne(context.TODO(), bson.M{"username": body.Username})
	if usernameExistsWithUsername.Err() == nil { // if no error, it means the user exists
		fmt.Print(usernameExistsWithUsername)
		logrus.Error("You already have an account with this username. Please login")
		response.HandleResponse(c, http.StatusBadRequest, "You already have an account with this username. Please login", nil)
		return
	}

	userExistsWithEmail := database.UserCollection.FindOne(context.TODO(), bson.M{"email": body.Email})
	if userExistsWithEmail.Err() != nil { // User does not exists with this email
		// hash the password
		hashedPassword := utils.HashPassword(body.Password)

		// generate verification code
		verificationCode := utils.GenerateVerificationCode()

		// create the user
		_, err := models.CreateUser(&models.User{
			Name:             body.Name,
			Username:         "",
			Email:            body.Email,
			Password:         hashedPassword,
			VerificationCode: verificationCode,
		})
		if err != nil {
			logrus.Errorf("Error creating user: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Error creating user", nil)
			return
		}

		// send verification email through rabbitmq and also send username
		queue.StartProducer(queue.EmailVerificationPayload{
			Email:    body.Email,
			Username: body.Username,
			Code:     verificationCode,
		})

		response.HandleResponse(c, http.StatusCreated, "User registered successfully. Please verify your email", nil)
	} else {
		var decodedUser models.User
		if err := userExistsWithEmail.Decode(&decodedUser); err != nil {
			logrus.Errorf("Error decoding the user: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
			return
		}

		if !decodedUser.IsEmailVerified {
			if decodedUser.VerificationCodeExpiresAt.Unix() < utils.GetCurrentDateTime() {
				logrus.Error("verification code expired. Please resend the verification code")
				response.HandleResponse(c, http.StatusBadRequest, "verification code expired. Please resend the verification code", nil)
				return
			}

			logrus.Error("Please verify your email to activate your account")
			response.HandleResponse(c, http.StatusBadRequest, "Please verify your email to activate your account", nil)
			return
		}
	}
}

func Login(c *gin.Context) {

}

func Logout(c *gin.Context) {

}
