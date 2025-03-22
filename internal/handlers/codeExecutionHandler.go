package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	codeexecutor "github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/code-executor"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExecuteQuestion(c *gin.Context) {
	var body codeexecutor.ExecuteQuestion
	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Errorf("Invalid request body: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	err := utils.ValidateRequest(body)
	if err != nil {
		logrus.Errorf("Error validating the request body: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	objectId, err := primitive.ObjectIDFromHex(body.QuestionID)
	if err != nil {
		logrus.Errorf("Invalid question id: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid question id", nil)
		return
	}

	// check if question exists
	result := database.DBClient.Database(config.Config.DATABASE_NAME).Collection(constants.QUESTION_COLLECTION).FindOne(context.TODO(), bson.M{"_id": objectId})
	if result.Err() != nil {
		logrus.Errorf("Question not found: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	var question models.Question
	if err := result.Decode(&question); err != nil {
		logrus.Errorf("Error decoding the question: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// get the data from context
	decodeUser, err := utils.GetDecodedUserFromContext(c)
	if err != nil {
		logrus.Errorf("Error getting decoded user: ExecuteQuestion API: %v", err)
		response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}

	// execute the question
	if body.Type == constants.RUN_QUESTION || body.Type == constants.SUBMIT_QUESTION {
		message := "Question Run Successful!"

		var codeSnippet string
		for _, snippet := range question.CodeSnippets {
			if strings.ToLower(string(snippet.Language)) == body.Language {
				codeSnippet = strings.TrimSpace(snippet.Code)
				break
			}
		}

		// generate the code by replacing placeholders
		code := utils.GenerateCodeTemplate(question.TestCases, body.Language, codeSnippet, body.Code)

		// run the code for given language in its container
		_, err := services.ExecuteCodeInDocker(body.Language, code)
		if err != nil {
			logrus.Errorf("Error running the code: ExecuteQuestion API: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		// create an entry into the database
		if body.Type == constants.SUBMIT_QUESTION {
			message = "Question Submission Successful!"

			_, err := models.CreateSubmission(&models.QuestionSubmission{
				QuestionID: body.QuestionID,
				UserID: decodeUser.ID,
				Status: "success",
				CreatedAt: time.Now(),
			})
			if err != nil {
				logrus.Errorf("Error creating submission: ExecuteQuestion API: %v", err)
				response.HandleResponse(c, http.StatusInternalServerError, "Something went wrong", nil)
				return
			}
		}

		response.HandleResponse(c, http.StatusOK, message, nil)
		return
	} else {
		logrus.Errorf("Invalid execution operation: RunQuestionHandler API: %v", nil)
		response.HandleResponse(c, http.StatusBadRequest, "Invalid execution operation", nil)
		return
	}
}