package llms

import (
	"context"
	"fmt"
	"strconv"
)

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
	resp, err := GeminiGenerateContent(ctx, prompt)
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
