package handlers

import (
	"context"
	"net/http"
	"reflect"
	"strings"

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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	difficulty := c.Query("difficulty")
	category := c.Query("category")

	var filter interface{}

	if difficulty != "" && category != "" {
		filter = bson.M{
			"difficulty": difficulty,
			"tags": bson.M{
				"$elemMatch": bson.M{"$eq": category},
			},
		}
	} else if difficulty != "" {
		filter = bson.M{
			"difficulty": difficulty,
		}
	} else if category != "" {
		filter = bson.M{
			"tags": bson.M{
				"$elemMatch": bson.M{"$eq": category},
			},
		}
	} else {
		filter = bson.M{}
	}

	cursor, err := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).Find(context.TODO(), filter)
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

func GetQuestionById(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Invalid question id: GetQuestionById API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid question id", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Errorf("Question not found: GetQuestionById API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "Question not found", nil)
		return
	}

	var question models.Question
	if err := result.Decode(&question); err != nil {
		logrus.Errorf("Error decoding the question: GetQuestionById API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Question retrieved successfully", question)
}

func UpdateQuestion(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Invalid question id: UpdateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid question id", nil)
		return
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Errorf("Question not found: UpdateQuestion API: %v", result.Err())
		response.HandleResponse(c, http.StatusNotFound, "Question not found", nil)
		return
	}

	var questionToUpdate models.Question
	if err := result.Decode(&questionToUpdate); err != nil {
		logrus.Errorf("Error decoding the question: UpdateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	var question request.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&question); err != nil {
		logrus.Errorf("Invalid request body: UpdateQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if question.Title != "" && question.Title != questionToUpdate.Title {
		questionToUpdate.Title = question.Title

		// change the slug
		words := strings.Split(question.Title, " ")
		slug := strings.Join(words, "-")
		questionToUpdate.Slug = slug
	}

	if question.Description != "" && question.Description != questionToUpdate.Description {
		questionToUpdate.Description = question.Description
	}

	if question.Difficulty != "" && question.Difficulty != questionToUpdate.Difficulty {
		questionToUpdate.Difficulty = question.Difficulty
	}

	if question.Tags != nil && !reflect.DeepEqual(question.Tags, questionToUpdate.Tags) {
		questionToUpdate.Tags = question.Tags
	}

	if question.Companies != nil && !reflect.DeepEqual(question.Companies, questionToUpdate.Companies) {
		questionToUpdate.Companies = question.Companies
	}

	if question.Hints != nil && !reflect.DeepEqual(question.Hints, questionToUpdate.Hints) {
		questionToUpdate.Hints = question.Hints
	}

	if question.TestCases != nil && !reflect.DeepEqual(question.TestCases, questionToUpdate.TestCases) {
		questionToUpdate.TestCases = question.TestCases
	}

	if question.CodeSnippets != nil && !reflect.DeepEqual(question.CodeSnippets, questionToUpdate.CodeSnippets) {
		questionToUpdate.CodeSnippets = question.CodeSnippets
	}

	updateStage := bson.M{
		"$set": bson.M{
			"title":        questionToUpdate.Title,
			"description":  questionToUpdate.Description,
			"difficulty":   questionToUpdate.Difficulty,
			"tags":         questionToUpdate.Tags,
			"companies":    questionToUpdate.Companies,
			"hints":        questionToUpdate.Hints,
			"testCases":    questionToUpdate.TestCases,
			"codeSnippets": questionToUpdate.CodeSnippets,
			"slug":         questionToUpdate.Slug,
		},
	}
	res, updateErr := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).UpdateOne(context.TODO(), bson.M{"_id": objectId}, updateStage)
	if updateErr != nil {
		logrus.Errorf("Error updating question: UpdateQuestion API: %v", updateErr)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	if res.ModifiedCount == 0 {
		response.HandleResponse(c, http.StatusOK, "No changes made", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Question updated successfully", questionToUpdate)
}

func DeleteQuestion(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("Invalid question id: DeleteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid question id", nil)
	}

	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).FindOneAndDelete(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Error("Question not found: DeleteQuestion API")
		response.HandleResponse(c, http.StatusNotFound, "Question not found", nil)
		return
	}

	response.HandleResponse(c, http.StatusOK, "Question deleted successfully", nil)
}

func SearchQuestions(c *gin.Context) {}

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
