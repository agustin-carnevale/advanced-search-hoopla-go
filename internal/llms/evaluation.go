package llms

import (
	"context"
	"encoding/json"
	"fmt"
)

func EvaluateResults(ctx context.Context, query string, resultsListStr string) ([]int, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Rate how relevant each result is to this query on a 0-3 scale:
    Query: "%s"

    Results:
    %s

    Scale:
    - 3: Highly relevant
    - 2: Relevant
    - 1: Marginally relevant
    - 0: Not relevant

    Do NOT give any numbers out than 0, 1, 2, or 3.

    Return ONLY the scores in the same order you were given the documents. Return a valid JSON list, nothing else. For example:

    [2, 0, 3, 2, 0, 1]
		`,
		query,
		resultsListStr,
	)

	// Call to llm
	jsonData, _, err := GeminiGenerateContent(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var scoresList []int
	err = json.Unmarshal([]byte(jsonData), &scoresList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list of eval scores from response %q: %w", jsonData, err)
	}

	return scoresList, nil
}
