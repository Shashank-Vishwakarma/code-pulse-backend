package challenge

type ChallengeRequest struct {
	Topic string `json:"topic" validate:"required"`
}