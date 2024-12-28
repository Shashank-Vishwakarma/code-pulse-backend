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
	AUTH_API_FORGOT_PASSWORD_ENDPOINT          = "/forgot-password"
	AUTH_API_RESEND_VERIFICATION_CODE_ENDPOINT = "/:username/resend-verification-code"

	// Question API Endpoints
	QUESTION_API_BASE_ENDPOINT              = "/api/v1/questions"
	QUESTION_API_CREATE_ENDPOINT            = "/create"
	QUESTION_API_GET_ALL_ENDPOINT           = "/"
	QUESTION_API_GET_BY_ID_ENDPOINT         = "/:id"
	QUESTION_API_UPDATE_ENDPOINT            = "/:id"
	QUESTION_API_DELETE_ENDPOINT            = "/:id"
	QUESTION_API_SEARCH_ENDPOINT            = "/search"
	QUESTION_API_GET_BY_CATEGORY_ENDPOINT   = "/category/:category"
	QUESTION_API_GET_BY_DIFFICULTY_ENDPOINT = "/difficulty/:difficulty"

	// Blog API Endpoints
	BLOG_API_BASE_ENDPOINT           = "/api/v1/blogs"
	BLOG_API_CREATE_ENDPOINT         = "/create"
	BLOG_API_GET_ALL_ENDPOINT        = "/"
	BLOG_API_GET_BY_ID_ENDPOINT      = "/:id"
	BLOG_API_UPDATE_ENDPOINT         = "/:id"
	BLOG_API_DELETE_ENDPOINT         = "/:id"
	BLOG_API_SEARCH_ENDPOINT         = "/search"
	BLOG_API_GET_BY_USER_ID_ENDPOINT = "/user/:id"
)
