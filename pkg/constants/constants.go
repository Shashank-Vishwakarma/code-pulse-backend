package constants

const (
	// code execution
	RUN_QUESTION    = "run"
	SUBMIT_QUESTION = "submit"

	// Database
	USER_COLLECTION            = "users"
	QUESTION_COLLECTION        = "questions"
	BLOG_COLLECTION            = "blogs"
	COMMENT_COLLECTION         = "comments"
	CODE_SUBMISSION_COLLECTION = "submissions"
	CHALLENGE_COLLECTION       = "challenges"

	// Auth API Endpoints
	AUTH_API_BASE_ENDPOINT                     = "/api/v1/auth"
	AUTH_API_LOGIN_ENDPOINT                    = "/login"
	AUTH_API_REGISTER_ENDPOINT                 = "/register"
	AUTH_API_LOGOUT_ENDPOINT                   = "/logout"
	AUTH_API_EMAIL_VERIFY_ENDPOINT             = "/email/verify"
	AUTH_API_FORGOT_PASSWORD_ENDPOINT          = "/forgot-password"
	AUTH_API_RESEND_VERIFICATION_CODE_ENDPOINT = "/:username/resend-verification-code"

	// Question API Endpoints
	QUESTION_API_BASE_ENDPOINT                            = "/api/v1/questions"
	QUESTION_API_CREATE_ENDPOINT                          = "/create"
	QUESTION_API_GET_ALL_ENDPOINT                         = "/"
	QUESTION_API_GET_BY_ID_ENDPOINT                       = "/:id"
	QUESTION_API_UPDATE_ENDPOINT                          = "/:id"
	QUESTION_API_DELETE_ENDPOINT                          = "/:id"
	QUESTION_API_GET_BY_USER_ENDPOINT                     = "/user"
	QUESTION_API_GET_QUESTIONS_SUBMITTED_BY_USER_ENDPOINT = "/submitted"

	// Blog API Endpoints
	BLOG_API_BASE_ENDPOINT           = "/api/v1/blogs"
	BLOG_API_CREATE_ENDPOINT         = "/create"
	BLOG_API_GET_ALL_ENDPOINT        = "/"
	BLOG_API_GET_BY_ID_ENDPOINT      = "/:id"
	BLOG_API_UPDATE_ENDPOINT         = "/:id"
	BLOG_API_DELETE_ENDPOINT         = "/:id"
	BLOG_API_GET_BY_USER_ID_ENDPOINT = "/user"
	BLOG_API_CREATE_COMMENT_ENDPOINT = "/:id/comments"

	// Code Execution Endpoints
	CODE_EXECUTION_API_BASE_ENDPOINT = "/api/v1/questions/:id/execute/"

	// Quiz API Endpoints
	CHALLENGE_API_BASE_ENDPOINT                          = "/api/v1/challenges"
	CHALLENGE_API_ALL_CHALLENGES_ENDPOINT                = "/all"
	CHALLENGE_API_CREATE_ENDPOINT                        = "/create"
	CHALLENGE_API_GET_BY_ID_ENDPOINT                     = "/:id"
	CHALLENGE_API_DELETE_ENDPOINT                        = "/:id"
	CHALLENGE_API_GET_ALL_BY_USER_ID_ENDPOINT            = "/user"
	CHALLENGE_API_SUBMIT_CHALLENGE_ENDPOINT              = "/:id/submit"
	CHALLENGE_API_GET_CORRECT_ANSWERS_CHALLENGE_ENDPOINT = "/:id/answers"
	CHALLENGE_API_GET_CHALLENGES_TAKEN_ENDPOINT          = "/taken"
)
