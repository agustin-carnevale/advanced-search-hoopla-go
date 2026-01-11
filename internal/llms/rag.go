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

func ResultsWithCitations(ctx context.Context, query string, resultsStr string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Answer the question or provide information based on the provided documents.
			This should be tailored to Hoopla users. Hoopla is a movie streaming service.

			If not enough information is available to give a good answer, say so but give as good of an answer as you can while citing the sources you have.

			Query: %s

			Documents:
			%s

			Instructions:
			- Provide a comprehensive answer that addresses the query
			- Cite sources using [1], [2], etc. format when referencing information
			- If sources disagree, mention the different viewpoints
			- If the answer isn't in the documents, say "I don't have enough information"
			- Be direct and informative

			Answer:
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

func AnswerQuestionFromResults(ctx context.Context, query string, resultsStr string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Answer the following question based on the provided documents.

			Question: %s

			Documents:
			%s

			General instructions:
			- Answer directly and concisely
			- Use only information from the documents
			- If the answer isn't in the documents, say "I don't have enough information"
			- Cite sources when possible

			Guidance on types of questions:
			- Factual questions: Provide a direct answer
			- Analytical questions: Compare and contrast information from the documents
			- Opinion-based questions: Acknowledge subjectivity and provide a balanced view

			Answer:
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
