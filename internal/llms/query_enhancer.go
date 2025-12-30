package llms

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func QueryEnhanceSpell(ctx context.Context, query string) (string, error) {
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

	// Build the text prompt
	prompt := fmt.Sprintf(`Fix any spelling errors in this movie search query.
		Only correct obvious typos. Don't change correctly spelled words.

		Query: "%s"

		If no errors, return the original query.
		If the query was corrected, just return the new query (ready to use, without additional text).`,
		query,
	)

	// Wrap into a *genai.Content slice
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

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
