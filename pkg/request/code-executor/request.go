package codeexecutor

type ExecuteQuestion struct {
	QuestionID string `json:"question_id" validate:"required"`
	Language   string `json:"language" validate:"required"`
	Code       string `json:"code" validate:"required"`
	Type       string `json:"type" validate:"required"` // "run" or "submit"
}
