package utils

import (
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/challenge"
)

// Each question carried 1 marks
const MAX_SCORE = 10

func CalculateChallengeScore(answers []challenge.ChallengeSubmittedQuestionAnswer, data []models.ChallegeQuestion) string {
	var score int
	for i, answer := range answers {
		if answer.Answer== data[i].CorrectAnswer {
			score += 1
		}
	}

	return fmt.Sprintf("%d/%d", score, MAX_SCORE)
}