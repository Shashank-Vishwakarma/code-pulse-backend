package routes

import (
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func CodeExecutionRoutes(r *gin.Engine) {
	codeExecutionRouteGroup := r.Group(constants.CODE_EXECUTION_API_BASE_ENDPOINT)

	codeExecutionRouteGroup.Use(middlewares.Authorization())

	codeExecutionRouteGroup.POST("/", handlers.ExecuteQuestion)
}