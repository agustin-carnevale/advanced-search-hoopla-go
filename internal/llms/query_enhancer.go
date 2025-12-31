package llms

import (
	"context"
	"fmt"
)

func QueryEnhanceSpell(ctx context.Context, query string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Fix any spelling errors in this movie search query.
		Only correct obvious typos. Don't change correctly spelled words.

		Query: "%s"

		If no errors, return the original query.
		If the query was corrected, just return the new query (ready to use, without additional text).`,
		query,
	)

	// Call llm
	return GeminiGenerateContent(ctx, prompt)
}

func QueryEnhanceRewrite(ctx context.Context, query string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Rewrite this movie search query to be more specific and searchable.
	
		Original: "%s"

    Consider:
    - Common movie knowledge (famous actors, popular films)
    - Genre conventions (horror = scary, animation = cartoon)
    - Keep it concise (under 10 words)
    - It should be a google style search query that's very specific
    - Don't use boolean logic

    Examples:

    - "that bear movie where leo gets attacked" -> "The Revenant Leonardo DiCaprio bear attack"
    - "movie about bear in london with marmalade" -> "Paddington London marmalade"
    - "scary movie with bear from few years ago" -> "bear horror movie 2015-2020"
    
    If you don't see any possible improvement, return the original query.
    If the query was optimized, just return the new query (ready to use, without additional text).`,
		query,
	)

	// Call llm
	return GeminiGenerateContent(ctx, prompt)
}
func QueryEnhanceExpand(ctx context.Context, query string) (string, error) {

	// Build the text prompt
	prompt := fmt.Sprintf(`Expand this movie search query with related terms.
		
		Add synonyms and related concepts that might appear in movie descriptions.
    Keep expansions relevant and focused.
    This will be appended to the original query.

    Examples:

    - "scary bear movie" -> "scary horror grizzly bear movie terrifying film"
    - "action movie with bear" -> "action thriller bear chase fight adventure"
    - "comedy with bear" -> "comedy funny bear humor lighthearted"
    

		Query: "%s"

 		If you don't see any possible improvement, return the original query.
    If the query was optimized, just return the new query (ready to use, without additional text).`,
		query,
	)

	// Call llm
	return GeminiGenerateContent(ctx, prompt)
}

func PreProcessQuery(ctx context.Context, query string, enhance string) (string, error) {
	switch enhance {
	case "spell":
		return QueryEnhanceSpell(ctx, query)
	case "rewrite":
		return QueryEnhanceRewrite(ctx, query)
	case "expand":
		return QueryEnhanceExpand(ctx, query)
	default:
		return query, nil
	}
}
