package utils

import (
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/request/challenge"
)

// Each question carried 1 marks
const MAX_SCORE = 10

func CalculateChallengeScore(answers []challenge.ChallengeSubmittedQuestionAnswer, data []models.ChallegeQuestion) string {
	score := 0

	for _, answer := range answers {
		for _, question := range data {
			if (answer.Question == question.Question) && (answer.Answer == question.CorrectAnswer) {
				score += 1
			}
		}
	}

	return fmt.Sprintf("%d/%d", score, MAX_SCORE)
}