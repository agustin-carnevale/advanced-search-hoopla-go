package llms

import (
	"context"
	"fmt"
)

func RewriteQueryFromImage(ctx context.Context, query string, img []byte, mime string) (string, int, error) {
	// Build the system prompt
	prompt := fmt.Sprintf(`Given the included image and text query, rewrite the text query to improve search results from a movie database. 
Make sure to:
- Synthesize visual and textual information
- Focus on movie-specific details (actors, scenes, style, etc.)
- Return only the rewritten query, without any additional commentary

Text query: %s`, query)

	// Call the generic multimodal generation function
	text, totalTokens, err := GeminiMultimodalGenerateContent(ctx, prompt, img, mime)
	if err != nil {
		return "", 0, err
	}

	return text, totalTokens, nil
}
