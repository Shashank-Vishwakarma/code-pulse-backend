package routes

import (
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func QuestionRoutes(r *gin.Engine) {
	questionRouteGroup := r.Group(constants.QUESTION_API_BASE_ENDPOINT)

	questionRouteGroup.Use(middlewares.Authorization())

	questionRouteGroup.POST(constants.QUESTION_API_CREATE_ENDPOINT, middlewares.RateLimiter(5, time.Hour), handlers.CreateQuestion)
	questionRouteGroup.GET(constants.QUESTION_API_GET_ALL_ENDPOINT, handlers.GetAllQuestions)
	questionRouteGroup.GET(constants.QUESTION_API_GET_BY_ID_ENDPOINT, handlers.GetQuestionById)
	questionRouteGroup.PUT(constants.QUESTION_API_UPDATE_ENDPOINT, handlers.UpdateQuestion)
	questionRouteGroup.DELETE(constants.QUESTION_API_DELETE_ENDPOINT, handlers.DeleteQuestion)
	questionRouteGroup.GET(constants.QUESTION_API_GET_BY_USER_ENDPOINT, handlers.GetQuestionsByUser)
}
