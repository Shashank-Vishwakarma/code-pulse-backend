package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/sirupsen/logrus"
)

var prompt = `
Generate 10 multiple-choice questions on the topic: %s. 

Each question should include:
1. A clear and concise question statement.
2. Four answer options.
3. The correct answer to each question.

Make the questions suitable for a quiz format, engaging, and of moderate difficulty. 
Ensure the questions are relevant to the topic and avoid ambiguity.

Please strictly provide a json response with the following format (Below is an example response that you should provide):
{
	"topic": "sql",
	"questions": [
		{
			"question": "question statement",
			"options": [option1, option2, option3, option4],
			"correct_answer": "option2"
		},
		{
			"question": "question text",
			"options": [option1, option2, option3, option4],
			"correct_answer": "option4"
		},
		... and so on for all 10 questions.
	]
}

This response should be a valid json.
`

const url = "https://api.groq.com/openai/v1/chat/completions"

type AIResponse struct {
	Questions []models.ChallegeQuestion `json:"questions"`
}

type Root struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int64       `json:"created"`
	Model             string      `json:"model"`
	Choices           []Choice    `json:"choices"`
	Usage             Usage       `json:"usage"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

func GenerateAIResponse(topic string) (AIResponse, error) {
	prompt = fmt.Sprintf(prompt, topic)

	payload := map[string]interface{}{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logrus.Printf("Error marshaling JSON: %v", err)
		return AIResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		logrus.Printf("Error creating request: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Config.GROQ_API_KEY))

	// make the request
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logrus.Printf("Could not make request: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}
	defer response.Body.Close()

	// Read the response body
	b, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Printf("Error reading response body: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}

	// Unmarshal the response body into AIResponse type
	var result Root
	err = json.Unmarshal(b, &result)
	if err != nil {
		logrus.Printf("Error unmarshalling the AI response: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}

	content := ""
	if result.Choices[0].Message.Content[0] == '`' {
		content = strings.ReplaceAll(result.Choices[0].Message.Content, "`", "")
	}

	// create AIResponse
	var res AIResponse
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		logrus.Printf("Error unmarshalling the content: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}

	return res, nil
}