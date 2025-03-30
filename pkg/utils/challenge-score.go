package utils

import (
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
)

// Each question carried 1 marks
const MAX_SCORE = 10

func CalculateChallengeScore(answers []map[string]string, data []models.ChallegeQuestion) string {
	var score int
	for i, answer := range answers {
		if answer["answer"] == data[i].CorrectAnswer {
			score += 1
		}
	}

	return fmt.Sprintf("%d/%d", score, MAX_SCORE)
}