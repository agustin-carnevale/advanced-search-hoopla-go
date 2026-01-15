package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type CohereRerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n,omitempty"`
}

type CohereRerankResponse struct {
	Results []struct {
		Index          int     `json:"index"`
		RelevanceScore float64 `json:"relevance_score"`
	} `json:"results"`
}

func ReRankDoc(ctx context.Context, query string, title string, description string) (float64, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`You are a movie search relevance rater. Rate how well this movie matches the search query on a scale of 0-10.

Query: "%s"
Movie Title: %s
Movie Description: %s

Scoring Guidelines:
- 10: Perfect match - directly matches all aspects of the query
- 8-9: Very relevant - matches most aspects, highly appropriate
- 6-7: Somewhat relevant - matches some aspects but not all
- 4-5: Partially relevant - matches a few aspects but misses key elements
- 2-3: Slightly relevant - only tangentially related
- 0-1: Not relevant - does not match the query

Consider:
- Direct relevance to query terms and concepts
- User intent (what they're looking for)
- Content appropriateness for the query context

IMPORTANT: Use the FULL 0-10 scale. Do NOT just return 0 or 1. Provide nuanced scores based on relevance.

Examples:
- Query: "family movie about bears" | Movie: "Paddington" | Score: 9
- Query: "family movie about bears" | Movie: "The Revenant" | Score: 2
- Query: "family movie about bears" | Movie: "Horror movie" | Score: 0

Return ONLY a single integer from 0-10, nothing else. No explanation, no text, just the number.`,
		query,
		title,
		description,
	)

	// Call to llm
	resp, _, err := GeminiGenerateContent(ctx, prompt)
	if err != nil {
		return 0.0, err
	}
	if resp == "" {
		// If LLM returned empty string, just treat as 0 score
		return 0.0, nil
	}

	score, err := strconv.Atoi(resp)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse score from response %q: %w", resp, err)
	}

	// Just for debugging purposes:
	// fmt.Printf("ReRanking Movie: %s\n", title)
	// fmt.Printf("LLM resp: %s\n", resp)
	// fmt.Printf("Score: %d\n", score)
	// fmt.Printf("Float Score: %.3f\n\n", float64(score))

	return float64(score), nil
}

func ReRankDocsBatch(ctx context.Context, query string, docsListStr string) ([]int, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Rank these movies by relevance to the search query.

      Query: "%s"
      Movies: 
      %s

      Consider:
      - Direct relevance to query
      - User intent (what they're looking for)
      - Content appropriateness

    	Return ONLY the IDs in order of relevance (best match first). Return a valid JSON list, nothing else. For example:

    	[75, 12, 34, 2, 1]
		`,
		query,
		docsListStr,
	)

	// Call to llm
	jsonData, _, err := GeminiGenerateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var rankedIdsList []int
	err = json.Unmarshal([]byte(jsonData), &rankedIdsList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list of ids from response %q: %w", jsonData, err)
	}

	return rankedIdsList, nil
}

// Using Cohere as the most similar approach to directly using/implementing a CrossEncoder
// Why this is a cross-encoder replacement:
//   - Cohere rerank jointly encodes query + document
//   - Uses a fine-tuned relevance model
//   - Returns comparable scalar relevance scores
//   - Deterministic & optimized for ranking
//
// Note: This is not prompt-based ranking.
func CohereRerankCrossEncoder(
	ctx context.Context,
	query string,
	docs []string,
) ([]float64, error) {
	apiKey := os.Getenv("COHERE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("COHERE_API_KEY not set")
	}

	reqBody := CohereRerankRequest{
		Model:     "rerank-english-v3.0",
		Query:     query,
		Documents: docs,
	}

	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.cohere.ai/v1/rerank",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cohere rerank failed: %s", string(b))
	}

	var out CohereRerankResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	// map scores back to original order
	scores := make([]float64, len(docs))
	for _, r := range out.Results {
		scores[r.Index] = r.RelevanceScore
	}

	return scores, nil
}
