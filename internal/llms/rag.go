package llms

import (
	"context"
	"fmt"
)

func ResultsAugmentation(ctx context.Context, query string, resultsStr string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Answer the question or provide information based on the provided documents. 
			This should be tailored to Hoopla users. Hoopla is a movie streaming service.
			Query: %s

			Documents:
			%s

			Provide a comprehensive answer that addresses the query:
		`,
		query,
		resultsStr,
	)

	// Call to llm
	response, err := GeminiGenerateContent(ctx, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

func SummarizeResults(ctx context.Context, query string, resultsStr string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Provide information useful to this query by synthesizing information from multiple search results in detail.
      The goal is to provide comprehensive information so that users know what their options are.
      Your response should be information-dense and concise, with several key pieces of information about the genre, plot, etc. of each movie.
      This should be tailored to Hoopla users. Hoopla is a movie streaming service.

      Query: %s

      Search Results:
      %s

      Provide a comprehensive 3-4 sentence answer that combines information from multiple sources:
		`,
		query,
		resultsStr,
	)

	// Call to llm
	response, err := GeminiGenerateContent(ctx, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
