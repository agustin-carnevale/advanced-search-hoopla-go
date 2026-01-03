package llms

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func GeminiGenerateContent(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	// Wrap into a *genai.Content slice
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	// TODO: fix this and make it work
	// genConfig := &genai.GenerateContentConfig{
	// 	MaxOutputTokens: 8,
	// }

	// Call GenerateContent
	response, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", contents, nil)
	if err != nil {
		return "", fmt.Errorf("generate content error: %w", err)
	}

	// Extract text
	if len(response.Candidates) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", nil
}
