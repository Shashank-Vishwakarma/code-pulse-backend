package request

import "github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"

// Auth requests
type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Username        string `json:"username" validate:"required,min=6,max=20"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=20"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=20"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8,max=20"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,min=6,max=6"`
}

type ForgotPasswordRequest struct {
	Username        string `json:"username" validate:"required,min=6,max=20"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=20"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=20"`
}

// Question requests
type CreateQuestionRequest struct {
	Title        string               `json:"title" validate:"required,min=5"`
	Description  string               `json:"description" validate:"required"`
	Difficulty   models.Difficulty    `json:"difficulty" validate:"required"`
	Tags         []string             `json:"tags" validate:"required"`
	Companies    []string             `json:"companies"`
	Hints        []string             `json:"hints"`
	TestCases    []models.TestCase    `json:"testCases" validate:"required"`
	CodeSnippets []models.CodeSnippet `json:"codeSnippets" validate:"required"`
}

type UpdateQuestionRequest struct {
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Difficulty   models.Difficulty    `json:"difficulty"`
	Tags         []string             `json:"tags"`
	Companies    []string             `json:"companies"`
	Hints        []string             `json:"hints"`
	TestCases    []models.TestCase    `json:"testCases"`
	CodeSnippets []models.CodeSnippet `json:"codeSnippets"`
}

// Blog requests
type CreateBlogRequest struct {
	Title           string `form:"title" validate:"required,min=5"`
	Body            string `form:"body" validate:"required"`
	IsBlogPublished bool   `form:"isBlogPublished"`
}

type UpdateBlogRequest struct {
	Title           string `json:"title"`
	Body            string `json:"body"`
	IsBlogPublished bool   `json:"isBlogPublished"`
}
