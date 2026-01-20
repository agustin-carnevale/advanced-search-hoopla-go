package llms

import (
	"context"
	"fmt"
	"os"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
	"google.golang.org/genai"
)

func MultimodalEmbedText(ctx context.Context, text string) ([]float32, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: text}}},
	}

	resp, err := client.Models.EmbedContent(ctx, "models/embedding-001", contents, nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("Embedding was not generated correctly.")
	}

	return resp.Embeddings[0].Values, nil
}

func BuildMultimodalEmbeddingsBatch(
	ctx context.Context,
	documents []model.Movie,
) ([][]float32, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	contents := make([]*genai.Content, 0, len(documents))
	for _, doc := range documents {
		text := doc.Title + ": " + doc.Description
		contents = append(contents, &genai.Content{
			Parts: []*genai.Part{{Text: text}},
		})
	}

	var allEmbeddings [][]float32
	batchSize := 100

	for i := 0; i < len(contents); i += batchSize {
		end := i + batchSize
		if end > len(contents) {
			end = len(contents)
		}

		batch := contents[i:end]
		resp, err := client.Models.EmbedContent(ctx, "models/embedding-001", batch, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embeddings for batch %d-%d: %w", i, end, err)
		}

		if resp == nil || len(resp.Embeddings) != len(batch) {
			if resp != nil && len(resp.Embeddings) > 0 {
				if len(resp.Embeddings) != len(batch) {
					return nil, fmt.Errorf("mismatch in batch %d-%d: expected %d embeddings, got %d", i, end, len(batch), len(resp.Embeddings))
				}
			}
		}

		for _, emb := range resp.Embeddings {
			allEmbeddings = append(allEmbeddings, emb.Values)
		}
	}

	return allEmbeddings, nil
}

// WORKAROUND: The 'models/multimodal-embedding-001' model is currently unavailable (returning 404).
// To enable image search, we use a "Describe-then-Embed" strategy:
// 1. Use a generative model (gemini-3-flash-preview) to describe the image in text.
// 2. Embed that text description using the standard text embedding model (embedding-001).
// This allows us to search against our existing movie text embeddings.
func EmbedImage(
	ctx context.Context,
	imageBytes []byte,
	mime string,
) ([]float32, error) {

	// Describe the image using a generative model
	prompt := "Describe this image in detail, focusing on known movies, actors, mood, and setting, for the purpose of matching it with a movie in a database."
	description, _, err := GeminiMultimodalGenerateContent(ctx, prompt, imageBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("failed to describe image: %w", err)
	}

	// Embed the description
	return MultimodalEmbedText(ctx, description)
}

// TODO: extract this logic to a helper function to avoid repeating
// apiKey := os.Getenv("GEMINI_API_KEY")
// if apiKey == "" {
// 	return nil, fmt.Errorf("GEMINI_API_KEY not set")
// }
// client, err := genai.NewClient(ctx, &genai.ClientConfig{
// 	APIKey:  apiKey,
// 	Backend: genai.BackendGeminiAPI,
// })
// if err != nil {
// 	return nil, fmt.Errorf("failed to create client: %w", err)
// }
