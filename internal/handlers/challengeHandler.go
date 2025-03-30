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
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionData struct {
	Question string   `bson:"question" json:"question"`
	Options  []string `bson:"options" json:"options"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChallengeData struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Title   string             `bson:"title" json:"title"`
	Topic   string             `bson:"topic" json:"topic"`
	Difficulty string             `bson:"difficulty" json:"difficulty"`
	Data    []QuestionData `bson:"data" json:"data"`
	Score string `bson:"score" json:"score"`
	UserID    primitive.ObjectID    `bson:"user_id" json:"user_id"`
	UserData  User           `bson:"user_data" json:"user_data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func CreateChallenge(c *gin.Context) {
	var body challenge.ChallengeRequest
	if err := c.ShouldBind(&body); err != nil {
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

	aiResponse, err := utils.GenerateAIResponse(body.Topic, body.Difficulty)
	if err != nil {
		logrus.Errorf("Error generating AI response: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Error converting user id to object id: CreateChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// save to db
	_, err = models.CreateChallenge(&models.Challenge{
		Title: body.Title,
		Topic: body.Topic,
		Difficulty: body.Difficulty,
		UserID: userObjectId,
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
	challengeId := c.Param("id")

	challengeObjectId, err := primitive.ObjectIDFromHex(challengeId)
	if err != nil {
		logrus.Errorf("Error getting challenge id from request: GetChallengeById API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid Challenge Id", nil)
		return
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"_id": challengeObjectId}}}, // Exclude specific user_id
		{{"$lookup", bson.M{
			"from":         "users",       // Name of the users collection
			"localField":   "user_id",     // Field in the challenges collection
			"foreignField": "_id",         // Field in the users collection
			"as":           "user_data",  // Output array field for user data
		}}},
		{{"$unwind", bson.M{"path": "$user_data", "preserveNullAndEmptyArrays": true}}}, // Flatten user_data array
		{{"$project", bson.M{
			"user_data.password": 0,   // Exclude password from user data
			"user_data._id":      0,   // Optional: Exclude MongoDB's _id field for user data
		}}},
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).Aggregate(context.TODO(), pipeline)
	if err != nil {
		logrus.Errorf("Error getting the challenges for id: %s: GetChallengeById API: %v",challengeId, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var challenges []ChallengeData
	err = cursor.All(context.TODO(), &challenges)
	if err != nil {
		logrus.Errorf("Error getting challenges for id: %s: GetChallengeById API: %v", challengeId, err)
		response.HandleResponse(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	if len(challenges) == 0 {
		logrus.Errorf("Challenge not found for id: %s: GetChallengeById API", challengeId)
		response.HandleResponse(c, http.StatusNotFound, "Challenge not found", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Success", challenges[0])
}

func DeleteChallenge(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Error getting challenge id from request: DeleteChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid Challenge Id", nil)
		return
	}
	
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: DeleteChallenge API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Error getting user id: DeleteChallenge API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).FindOneAndDelete(context.Background(), bson.M{"_id": objectId, "user_id": userObjectId})
	if result.Err() != nil {
		logrus.Errorf("Error deleting the challenge with id: %s: DeleteChallenge API: %v", id, result.Err().Error())
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

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Error getting user id: GetAllChallengesByUserId API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"user_id": userObjectId}}}, // Exclude specific user_id
		{{"$lookup", bson.M{
			"from":         "users",       // Name of the users collection
			"localField":   "user_id",     // Field in the challenges collection
			"foreignField": "_id",         // Field in the users collection
			"as":           "user_data",  // Output array field for user data
		}}},
		{{"$unwind", bson.M{"path": "$user_data", "preserveNullAndEmptyArrays": true}}}, // Flatten user_data array
		{{"$project", bson.M{
			"user_data.password": 0,   // Exclude password from user data
			"user_data._id":      0,   // Optional: Exclude MongoDB's _id field for user data
		}}},
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).Aggregate(context.TODO(), pipeline)
	if err != nil {
		logrus.Errorf("Error getting the challenges for all users execpt user id: %s: GetAllChallengesByUserId API: %v", decodeUser.ID, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var challenges []ChallengeData
	err = cursor.All(context.TODO(), &challenges)
	if err != nil {
		logrus.Errorf("Error getting challenges for user id: %s: GetAllChallengesByUserId API: %v", userId, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Fetched challenges successfully", challenges)
}

// get all challenges except for current user
func GetAllChallenges(c *gin.Context) {
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: GetAllChallenges API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(decodeUser.ID)
	if err != nil {
		logrus.Errorf("Error getting user id: GetAllChallenges API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"user_id": bson.M{"$ne": userObjectId}}}}, // Exclude specific user_id
		{{"$lookup", bson.M{
			"from":         "users",       // Name of the users collection
			"localField":   "user_id",     // Field in the challenges collection
			"foreignField": "_id",         // Field in the users collection
			"as":           "user_data",  // Output array field for user data
		}}},
		{{"$unwind", bson.M{"path": "$user_data", "preserveNullAndEmptyArrays": true}}}, // Flatten user_data array
		{{"$project", bson.M{
			"user_data.password": 0,   // Exclude password from user data
			"user_data._id":      0,   // Optional: Exclude MongoDB's _id field for user data
		}}},
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).Aggregate(context.TODO(), pipeline)
	if err != nil {
		logrus.Errorf("Error getting the challenges for all users execpt user id: %s: GetAllChallenges API: %v", decodeUser.ID, err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var challenges []ChallengeData
	err = cursor.All(context.TODO(), &challenges)
	if err != nil {
		logrus.Errorf("Error decoding challenges: GetAllChallenges API: %v", err)
		response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Fetched challenges successfully", challenges)
}

func SubmitChallenge(c *gin.Context) {
	id := c.Param("id")

	challengeObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Error converting the challenge id to object id: SubmitChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	var body challenge.SubmitChallengeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Error binding request body: SubmitChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err = utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: SubmitChallenge API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// check if challenge with this id exist
	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).FindOne(
		context.TODO(),
		bson.M{
			"_id": challengeObjectId,
		},
	)

	if result.Err() != nil {
		logrus.Errorf("Challenge not found: SubmitChallenge API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "Challenge not found", nil)
		return
	}

	var challenge models.Challenge
	err = result.Decode(&challenge)
	if err != nil {
		logrus.Errorf("Error decoding challenge: SubmitChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	score := utils.CalculateChallengeScore(body.Answers, challenge.Data)

	// update the challenge score
	_, err = database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.CHALLENGE_COLLECTION).UpdateOne(
		context.TODO(),
		bson.M{
			"_id": challengeObjectId,
		},
		bson.M{
			"$set": bson.M{
				"score": score,
				"user_selected_answers": body.Answers,
			},
		},
	)
	if err != nil {
		logrus.Errorf("Error updating challenge score: SubmitChallenge API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Challenge submitted successfully", score)
}