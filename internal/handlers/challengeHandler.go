package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/challenge"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateChallenge(c *gin.Context) {
	var body challenge.ChallengeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := utils.ValidateRequest(body); 
	if err != nil {
		logrus.Errorf("Error validating the request body: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	aiResponse, err := utils.GenerateAIResponse(body.Topic)
	if err != nil {
		logrus.Errorf("Error generating AI response: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// save to db
	_, err = models.CreateChallenge(&models.Challenge{
		Topic: body.Topic,
		UserID: decodeUser.ID,
		Data: aiResponse.Questions,
	})
	if err != nil {
		logrus.Errorf("Error inserting challenge in databse: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusCreated, "Challenge created successfully", nil)
}

func GetChallengeById(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Error getting challenge id from request: GetChallengeById API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid Challenge Id", nil)
		return
	}

	var challenge struct {
		ID      primitive.ObjectID `bson:"_id" json:"id"`
		Topic   string             `bson:"topic" json:"topic"`
		Data    []struct {
			Question string   `bson:"question" json:"question"`
			Options  []string `bson:"options" json:"options"`
			// Exclude correct_answer from JSON response
		} `bson:"data" json:"data"`
		UserID    string    `bson:"user_id" json:"user_id"`
		CreatedAt time.Time `bson:"created_at" json:"created_at"`
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).FindOne(context.Background(), bson.M{"_id": objectId})
	err = result.Decode(&challenge)
	if err != nil {
		logrus.Errorf("Error decoding the challenge: %s: GetAllChallengesByUserId API: %v", id, err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: GetAllChallengesByUserId API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if decodeUser.ID != challenge.UserID {
		logrus.Errorf("Unauthorized user: %s: GetAllChallengesByUserId API: %v", id, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}	

	response.HandleResponse(c, http.StatusOK, "Fetched challenge successfully", challenge)
}

func DeleteChallenge(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Error getting challenge id from request: GetChallengeById API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid Challenge Id", nil)
		return
	}
	
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: GetAllChallengesByUserId API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).FindOneAndDelete(context.Background(), bson.M{"_id": objectId, "user_id": decodeUser.ID})
	if result.Err() != nil {
		logrus.Errorf("Error deleting the challenge with id: %s: GetChallengeById API: %v", id, result.Err().Error())
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Challenge deleted successfully", nil)
}

func GetAllChallengesByUserId(c *gin.Context) {
	userId := c.Param("userId")

	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: GetAllChallengesByUserId API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if decodeUser.ID != userId {
		logrus.Errorf("User not authorized: GetAllChallengesByUserId API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).Find(context.TODO(), bson.M{"user_id": userId})
	if err != nil {
		logrus.Errorf("Error getting the challenges for user id: %s: GetAllChallengesByUserId API: %v", userId, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var challenges []models.Challenge
	err = cursor.All(context.TODO(), &challenges)
	if err != nil {
		logrus.Errorf("Error getting challenges for user id: %s: GetAllChallengesByUserId API: %v", userId, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var resData struct{
		Topics interface{} `json:"topics"`
	}

	topics := []map[string]string{}
	for _, challenge := range challenges {
		topics = append(topics, map[string]string{
			"id": challenge.ID,
			"topic": challenge.Topic,
		})
	}

	resData.Topics = topics

	response.HandleResponse(c, http.StatusOK, "Fetched challenges successfully", resData)
}
