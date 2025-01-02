package handlers

import (
	"context"
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateQuestion(c *gin.Context) {
	var body request.CreateQuestionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: CreateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: CreateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// get the data from context
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	result, err := models.CreateQuestion(&models.Question{
		Title:        body.Title,
		Description:  body.Description,
		Difficulty:   body.Difficulty,
		Tags:         body.Tags,
		Companies:    body.Companies,
		Hints:        body.Hints,
		TestCases:    body.TestCases,
		CodeSnippets: body.CodeSnippets,
		AuthorID:     decodeUser.ID,
	})
	if err != nil {
		logrus.Errorf("Error creating question: CreateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Failed to create question", nil)
		return
	}

	response.HandleResponse(c, http.StatusCreated, "Question created successfully", result)
}

func GetAllQuestions(c *gin.Context) {
	options := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).Find(context.TODO(), bson.M{}, options)
	if err != nil {
		logrus.Errorf("Error getting all questions: GetAllQuestions API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var result []models.Question
	if err := cursor.All(context.TODO(), &result); err != nil {
		logrus.Errorf("Error decoding the questions: GetAllQuestions API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Questions retrieved successfully", result)
}

func GetQuestionById(c *gin.Context) {}

func UpdateQuestion(c *gin.Context) {}

func DeleteQuestion(c *gin.Context) {}

func SearchQuestions(c *gin.Context) {}

func GetQuestionsByCategory(c *gin.Context) {}

func GetQuestionsByDifficulty(c *gin.Context) {}

func GetQuestionsByUser(c *gin.Context) {
	// get the data from context
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: CreateBlog API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	options := options.Find().SetSort(bson.M{"createdAt": -1})
	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).Find(context.TODO(), bson.M{"authorId": decodeUser.ID}, options)
	if err != nil {
		logrus.Errorf("Error getting all questions: GetAllQuestions API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var result []models.Question
	if err := cursor.All(context.TODO(), &result); err != nil {
		logrus.Errorf("Error decoding the questions: GetAllQuestions API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Questions retrieved successfully", result)
}
