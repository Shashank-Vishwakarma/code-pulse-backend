package handlers

import (
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	request "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/auth"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func GetAllQuestions(c *gin.Context) {}

func GetQuestionById(c *gin.Context) {}

func UpdateQuestion(c *gin.Context) {}

func DeleteQuestion(c *gin.Context) {}

func SearchQuestions(c *gin.Context) {}

func GetQuestionsByCategory(c *gin.Context) {}

func GetQuestionsByDifficulty(c *gin.Context) {}
