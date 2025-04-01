package routes

import (
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func ChallengeRoutes(r *gin.Engine) {
	challengeRouteGroup := r.Group(constants.CHALLENGE_API_BASE_ENDPOINT)

	challengeRouteGroup.Use(middlewares.Authorization())

	challengeRouteGroup.GET(constants.CHALLENGE_API_ALL_CHALLENGES_ENDPOINT, handlers.GetAllChallenges)
	challengeRouteGroup.GET(constants.CHALLENGE_API_GET_BY_ID_ENDPOINT, handlers.GetChallengeById)
	challengeRouteGroup.GET(constants.CHALLENGE_API_GET_ALL_BY_USER_ID_ENDPOINT, handlers.GetAllChallengesByUserId)
	challengeRouteGroup.GET(constants.CHALLENGE_API_GET_CORRECT_ANSWERS_CHALLENGE_ENDPOINT, handlers.GetCorrectAnswersForChallenge)
	challengeRouteGroup.POST(constants.CHALLENGE_API_CREATE_ENDPOINT, handlers.CreateChallenge)
	challengeRouteGroup.DELETE(constants.CHALLENGE_API_DELETE_ENDPOINT, handlers.DeleteChallenge)
	challengeRouteGroup.POST(constants.CHALLENGE_API_SUBMIT_CHALLENGE_ENDPOINT, handlers.SubmitChallenge)
}