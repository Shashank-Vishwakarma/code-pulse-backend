package challenge

type QuizRequest struct {
	Title string `json:"title" validate:"required"`
}