package routes

import (
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func CodeExecutionRoutes(r *gin.Engine) {
	r.POST(constants.CODE_EXECUTION_API_BASE_ENDPOINT, middlewares.Authorization(), middlewares.RateLimiter(5, time.Minute), handlers.ExecuteQuestion)
	r.POST(constants.COMPILER_CODE_EXECUTION_API_ENDPOINT, middlewares.Authorization(), middlewares.RateLimiter(5, time.Minute), handlers.ExecuteCompilerCode)
}