package routes

import (
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/handlers"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/middlewares"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group(constants.AUTH_API_BASE_ENDPOINT)
	authGroup.POST(constants.AUTH_API_REGISTER_ENDPOINT, handlers.Register)
	authGroup.POST(constants.AUTH_API_LOGIN_ENDPOINT, handlers.Login)
	authGroup.POST(constants.AUTH_API_LOGOUT_ENDPOINT, middlewares.Authorization(), handlers.Logout)
	authGroup.POST(constants.AUTH_API_EMAIL_VERIFY_ENDPOINT, middlewares.Authorization(), handlers.VerifyEmail)
	authGroup.POST(constants.AUTH_API_FORGOT_PASSWORD_ENDPOINT, handlers.ForgotPassword)
	authGroup.POST(constants.AUTH_API_RESEND_VERIFICATION_CODE_ENDPOINT, middlewares.Authorization(), handlers.ResendVerificationCodeViaEmail)
}
