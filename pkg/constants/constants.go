package constants

const (
	// Database
	USER_COLLECTION = "users"

	// Auth API Endpoints
	AUTH_API_BASE_ENDPOINT                     = "/api/v1/auth"
	AUTH_API_LOGIN_ENDPOINT                    = "/login"
	AUTH_API_REGISTER_ENDPOINT                 = "/register"
	AUTH_API_LOGOUT_ENDPOINT                   = "/logout"
	AUTH_API_EMAIL_VERIFY_ENDPOINT             = "/email/verify"
	AUTH_API_FORGOT_PASSWORD_ENDPOINT          = "/:username/forgot-password"
	AUTH_API_RESEND_VERIFICATION_CODE_ENDPOINT = "/:username/resend-verification-code"
)
