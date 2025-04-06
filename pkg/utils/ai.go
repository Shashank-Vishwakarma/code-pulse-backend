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
You are an AI that generates quiz questions.

Your task is to generate exactly 10 multiple-choice questions on the topic: "%s".
The difficulty level of the challenge should be: "%s".

⚠️ IMPORTANT:
- Your response MUST be a valid JSON object.
- Do NOT include any explanations, markdown, code blocks, or extra text — ONLY the JSON.
- All keys and string values must be enclosed in double quotes.
- Use proper JSON syntax, with commas, brackets, and colons placed correctly.

🧠 FORMAT TO FOLLOW (strictly return this structure):
{
	"topic": "your topic here",
	"questions": [
		{
		"question": "Your question text",
		"options": ["Option 1", "Option 2", "Option 3", "Option 4"],
		"correct_answer": "One of the above options"
		},
		...
		(Total 10 questions)
	]
}

✅ Example (shortened for reference):
{
	"topic": "SQL",
	"questions": [
		{
		"question": "Which SQL keyword is used to retrieve data?",
		"options": ["INSERT", "SELECT", "UPDATE", "DELETE"],
		"correct_answer": "SELECT"
		},
		...
	]
}

Return ONLY valid JSON and nothing else.
`

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

func GenerateAIResponse(topic, difficulty string) (AIResponse, error) {
	prompt = fmt.Sprintf(prompt, topic, difficulty)

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

	req, err := http.NewRequest("POST", config.Config.GROQ_CHAT_COMPLETION_ENDPOINT, bytes.NewBuffer([]byte(jsonData)))
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

	content := result.Choices[0].Message.Content
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	// create AIResponse
	var res AIResponse
	err = json.Unmarshal([]byte(content), &res)
	if err != nil {
		logrus.Printf("Error unmarshalling the content: GenerateAIResponse: %v", err)
		return AIResponse{}, err
	}

	return res, nil
}