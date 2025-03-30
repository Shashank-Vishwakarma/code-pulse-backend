package challenge

type ChallengeRequest struct {
	Title      string `form:"title" validate:"required,min=5"`
	Topic      string `form:"topic" validate:"required"`
	Difficulty string `form:"difficulty" validate:"required"`
}

type SubmitChallengeRequest struct {
	Answers []map[string]string `json:"answers" validate:"required"`
}