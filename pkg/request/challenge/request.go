package challenge

type ChallengeRequest struct {
	Title      string `form:"title" validate:"required,min=5"`
	Topic      string `form:"topic" validate:"required"`
	Difficulty string `form:"difficulty" validate:"required"`
}

type ChallengeSubmittedQuestionAnswer struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
}

type SubmitChallengeRequest struct {
	Answers []ChallengeSubmittedQuestionAnswer `json:"answers" validate:"required"`
}