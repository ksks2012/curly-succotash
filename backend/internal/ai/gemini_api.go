package ai

import (
	"context"
	"curly-succotash/backend/global"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"
)

// GeminiClient handles API calls to Google Gemini
type GeminiClient struct {
	client *genai.Client
	model  string
	config *genai.GenerateContentConfig
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv(global.AISetting.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %s", err)
	}

	return &GeminiClient{
		client: client,
		model:  global.AISetting.Model,
		config: &genai.GenerateContentConfig{
			Temperature:      &global.AISetting.Temperature,
			TopP:             &global.AISetting.TopP,
			TopK:             &global.AISetting.TopK,
			MaxOutputTokens:  global.AISetting.MaxOutputTokens,
			ResponseMIMEType: "application/json",
		},
	}, nil
}

// GenerateContent generates content using Gemini API
func (c *GeminiClient) GenerateContent(prompt string) (string, error) {
	ctx := context.Background()
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		c.config,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %s", err)
	}

	text := result.Text()
	time.Sleep(500 * time.Millisecond)

	return text, nil
}

// Close is a no-op because genai.Client does not require closing resources
func (c *GeminiClient) Close() error {
	return nil
}
