package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AIService struct {
	apiURL  string
	openRouterToken string
	client  *http.Client
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []OpenRouterMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
	MaxTokens   int                 `json:"max_tokens"`
}

type OpenRouterChoice struct {
	Message OpenRouterMessage `json:"message"`
}

type OpenRouterResponse struct {
	Choices []OpenRouterChoice `json:"choices"`
}

func NewAIService() *AIService {
	openRouterToken := os.Getenv("OPENROUTER_API_TOKEN")
	if openRouterToken == "" {
		panic("OPENROUTER_API_TOKEN environment variable is required")
	}

	return &AIService{
		apiURL:  "https://openrouter.ai/api/v1/chat/completions",
		openRouterToken: openRouterToken,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateContent sends a request to OpenRouter API and returns generated content
func (ai *AIService) GenerateContent(prompt string) (string, int, error) {
	// Prepare the request for OpenRouter API
	requestBody := OpenRouterRequest{
		Model: "openai/gpt-3.5-turbo", // Using a reliable model
		Messages: []OpenRouterMessage{
			{
				Role:    "system",
				Content: "You are a helpful AI assistant that generates blog content. Provide informative, engaging, and well-structured responses.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	// Convert to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", ai.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+ai.openRouterToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/your-repo") // Required by OpenRouter
	req.Header.Set("X-Title", "Blog API") // Optional but recommended

	// Send request
	resp, err := ai.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var openRouterResponse OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResponse); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(openRouterResponse.Choices) == 0 {
		return "", 0, fmt.Errorf("no response from AI model")
	}

	// Extract generated text
	generatedText := openRouterResponse.Choices[0].Message.Content

	// Estimate token count (rough approximation: 1 token ≈ 4 characters)
	tokenCount := len(generatedText) / 4

	return generatedText, tokenCount, nil
}

// GenerateBlogContent generates blog content based on user's request
func (ai *AIService) GenerateBlogContent(topic string) (string, int, error) {
	prompt := fmt.Sprintf("Write a blog post about: %s. Make it informative, engaging, and well-structured.", topic)
	return ai.GenerateContent(prompt)
}

// SuggestTopics generates topic suggestions for blog posts
func (ai *AIService) SuggestTopics(category string) (string, int, error) {
	prompt := fmt.Sprintf("Suggest 5 interesting blog topics about: %s. Format as a numbered list.", category)
	return ai.GenerateContent(prompt)
}

// EnhanceContent improves existing blog content
func (ai *AIService) EnhanceContent(content string) (string, int, error) {
	prompt := fmt.Sprintf("Improve and enhance this blog content, making it more engaging and professional: %s", content)
	return ai.GenerateContent(prompt)
}
