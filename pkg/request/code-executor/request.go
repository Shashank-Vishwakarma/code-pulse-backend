package codeexecutor

type ExecuteQuestion struct {
	Language string `json:"language" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Type     string `json:"type" validate:"required"` // "run" or "submit"
}

type ExecuteCompilerCode struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required"`
}