package llms

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func GeminiGenerateContent(ctx context.Context, prompt string) (string, int, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", 0, fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", 0, fmt.Errorf("failed to create client: %w", err)
	}

	// Wrap into a *genai.Content slice
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	// Call GenerateContent
	response, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", contents, nil)
	if err != nil {
		return "", 0, fmt.Errorf("generate content error: %w", err)
	}

	// Extract text and token count
	var text string
	var totalTokens int
	if len(response.Candidates) > 0 {
		text = response.Candidates[0].Content.Parts[0].Text
	}
	if response.UsageMetadata != nil {
		totalTokens = int(response.UsageMetadata.TotalTokenCount)
	}

	return text, totalTokens, nil
}

func GeminiMultimodalGenerateContent(ctx context.Context, textPrompt string, img []byte, mime string) (string, int, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", 0, fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", 0, fmt.Errorf("failed to create client: %w", err)
	}

	// Build the parts for the multimodal prompt
	parts := []*genai.Part{
		{Text: textPrompt},
		{
			InlineData: &genai.Blob{
				MIMEType: mime,
				Data:     img,
			},
		},
	}

	// Wrap into a *genai.Content slice
	contents := []*genai.Content{
		{Parts: parts},
	}

	// Call GenerateContent
	response, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", contents, nil)
	if err != nil {
		return "", 0, fmt.Errorf("generate content error: %w", err)
	}

	// Extract text and token count from response
	var text string
	var totalTokens int
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		text = response.Candidates[0].Content.Parts[0].Text
	}
	if response.UsageMetadata != nil {
		totalTokens = int(response.UsageMetadata.TotalTokenCount)
	}

	if text == "" {
		return "", totalTokens, fmt.Errorf("no response text found")
	}

	return text, totalTokens, nil
}
