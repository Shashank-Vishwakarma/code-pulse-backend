package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/queue"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(c *gin.Context) {
	var body request.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: Register API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// validate the request body
	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: Register API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error validating the request body", nil)
		return
	}

	if body.Password != body.ConfirmPassword {
		logrus.Error("Passwords do not match: Register API")
		response.HandleResponse(c, http.StatusBadRequest, "Passwords do not match", nil)
		return
	}

	usernameExistsWithUsername := database.UserCollection.FindOne(context.TODO(), bson.M{"username": body.Username})
	if usernameExistsWithUsername.Err() == nil { // if no error, it means the user exists
		logrus.Error("You already have an account with this username. Please login: Register API")
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
		user, err := models.CreateUser(&models.User{
			Name:             body.Name,
			Username:         "",
			Email:            body.Email,
			Password:         hashedPassword,
			VerificationCode: verificationCode,
			QuestionsSubmitted: []string{},
			ChallengesTaken: []string{},
			Stats: models.Stats{
				QuestionsSubmitted: 0,
				QuestionsCreated:   0,
				BlogsCreated:       0,
				ChallengesCreated:  0,
				ChallengesTaken:    0,
			},
		})
		if err != nil {
			logrus.Errorf("Error creating user: Register API: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Error creating user", nil)
			return
		}

		// send verification email through rabbitmq and also send username
		queue.StartProducer(queue.EmailVerificationPayload{
			Email:    body.Email,
			Username: body.Username,
			Code:     verificationCode,
		})

		ID := ""
		if oid, ok := user.InsertedID.(primitive.ObjectID); ok {
			ID = oid.Hex()
		} else {
			logrus.Errorf("InsertedID is not a valid ObjectID: Register API: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
			return
		}

		// set jwt token in cookie
		token, err := utils.GenerateToken(utils.JWTPayload{
			ID:       ID,
			Name:     body.Name,
			Email:    body.Email,
			Username: body.Username,
		})
		if err != nil {
			logrus.Errorf("Error generating the token: Login API: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
			return
		}

		c.SetCookie(config.Config.JWT_TOKEN_COOKIE, token, 24*60*60, "/", "localhost", false, true)

		responseData := struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
			Stats    models.Stats `json:"stats"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:       ID,
			Name:     body.Name,
			Email:    body.Email,
			Username: body.Username,
			Stats: models.Stats{
				QuestionsSubmitted: 0,
				QuestionsCreated:   0,
				BlogsCreated:       0,
				ChallengesCreated:  0,
				ChallengesTaken:    0,
			},
			CreatedAt: time.Now(),
		}
		response.HandleResponse(c, http.StatusCreated, "User registered successfully. Please verify your email", responseData)
	} else {
		var decodedUser models.User
		if err := userExistsWithEmail.Decode(&decodedUser); err != nil {
			logrus.Errorf("Error decoding the user: Register API: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
			return
		}

		if !decodedUser.IsEmailVerified {
			if decodedUser.VerificationCodeExpiresAt.Unix() < utils.GetCurrentDateTime() {
				logrus.Error("verification code expired. Please resend the verification code: Register API")
				response.HandleResponse(c, http.StatusBadRequest, "verification code expired. Please resend the verification code", nil)
				return
			}

			logrus.Error("Please verify your email to activate your account: Register API")
			response.HandleResponse(c, http.StatusBadRequest, "Please verify your email to activate your account", nil)
			return
		}
	}
}

func Login(c *gin.Context) {
	var body request.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: Login API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// validate the request body
	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: Login API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error validating the request body", nil)
		return
	}

	filter := bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "email", Value: body.Identifier}}, bson.D{{Key: "username", Value: body.Identifier}}}}}
	result := database.UserCollection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		logrus.Errorf("User not found: Login API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	var decodedUser models.User
	if err := result.Decode(&decodedUser); err != nil {
		logrus.Errorf("Error decoding the user: Login API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	if !utils.CheckPasswordHash(body.Password, decodedUser.Password) {
		logrus.Error("Invalid Password: Login API")
		response.HandleResponse(c, http.StatusUnauthorized, "Invalid Password", nil)
		return
	}

	if !decodedUser.IsEmailVerified {
		logrus.Error("Please verify your email to activate your account: Login API")
		response.HandleResponse(c, http.StatusBadRequest, "Please verify your email to activate your account", nil)
		return
	}

	// set jwt token in cookie
	token, err := utils.GenerateToken(utils.JWTPayload{
		ID:       decodedUser.ID,
		Name:     decodedUser.Name,
		Email:    decodedUser.Email,
		Username: decodedUser.Username,
	})
	if err != nil {
		logrus.Errorf("Error generating the token: Login API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	c.SetCookie(config.Config.JWT_TOKEN_COOKIE, token, 24*60*60, "/", "", false, true)

	responseData := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Stats    models.Stats `json:"stats"`
		CreatedAt time.Time `json:"created_at"`
	}{
		ID:       decodedUser.ID,
		Name:     decodedUser.Name,
		Email:    decodedUser.Email,
		Username: decodedUser.Username,
		Stats:    decodedUser.Stats,
		CreatedAt: decodedUser.CreatedAt,
	}
	response.HandleResponse(c, http.StatusOK, "Login successful", responseData)
}

func Logout(c *gin.Context) {
	c.SetCookie(config.Config.JWT_TOKEN_COOKIE, "", -1, "/", "", false, true)
	response.HandleResponse(c, http.StatusOK, "Logout successful", nil)
}

func VerifyEmail(c *gin.Context) {
	var body request.VerifyEmailRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: VerifyEmail API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// validate the request body
	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: VerifyEmail API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error validating the request body", nil)
		return
	}

	result := database.UserCollection.FindOne(context.TODO(), bson.M{"email": body.Email})
	if result.Err() != nil {
		logrus.Errorf("User not found: VerifyEmail API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	var user models.User
	if err := result.Decode(&user); err != nil {
		logrus.Errorf("Error decoding the user: VerifyEmail API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	if user.IsEmailVerified {
		logrus.Error("Email already verified: VerifyEmail API")
		response.HandleResponse(c, http.StatusBadRequest, "Email is already verified. Please login", nil)
		return
	}

	if user.VerificationCode != body.Code {
		logrus.Error("Invalid verification code: VerifyEmail API")
		response.HandleResponse(c, http.StatusBadRequest, "Invalid verification code", nil)
		return
	}

	if user.VerificationCodeExpiresAt.Unix() < utils.GetCurrentDateTime() {
		logrus.Error("Verification code expired: VerifyEmail API")
		response.HandleResponse(c, http.StatusBadRequest, "Verification code expired. Please resend the verification code", nil)
		return
	}

	// get the data from context
	decodedUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: VerifyEmail API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	result = database.UserCollection.FindOneAndUpdate(context.TODO(), bson.M{"email": body.Email}, bson.M{
		"$set": bson.M{
			"username":                     decodedUser.Username,
			"is_email_verified":            true,
			"verification_code":            "",
			"verification_code_expires_at": time.Time{},
		},
	})
	if result.Err() != nil {
		logrus.Errorf("Error updating the user: VerifyEmail API: %v", result.Err())
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Email Verified sucessfully", nil)
}

func ForgotPassword(c *gin.Context) {
	var body request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: ForgotPassword API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// validate the request body
	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: ForgotPassword API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Error validating the request body", nil)
		return
	}

	if body.Password != body.ConfirmPassword {
		logrus.Error("Passwords do not match: ForgotPassword API")
		response.HandleResponse(c, http.StatusBadRequest, "Passwords do not match", nil)
		return
	}

	result := database.UserCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"$and": bson.A{
				bson.M{"email": body.Email},
				bson.M{"username": body.Username},
			},
		},
		bson.M{
			"$set": bson.M{
				"password": utils.HashPassword(body.Password),
			},
		},
	)
	if result.Err() != nil {
		logrus.Errorf("Error updating the user: ForgotPassword API: %v", result.Err())
		response.HandleResponse(c, http.StatusInternalServerError, "User not found with the given email and username", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Password reset successful", nil)
}

func ResendVerificationCodeViaEmail(c *gin.Context) {

}
