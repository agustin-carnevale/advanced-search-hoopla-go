package tokenizer

import (
	"strings"
	"unicode"

	"github.com/reiver/go-porterstemmer"
)

// HasMatchingToken checks if any token from queryTokens is contained within any token of titleTokens.
func HasMatchingToken(queryTokens, titleTokens []string) bool {
	for _, q := range queryTokens {
		for _, t := range titleTokens {
			if strings.Contains(t, q) {
				return true
			}
		}
	}
	return false
}

// preprocessText converts text to lowercase and removes punctuation.
func preprocessText(text string) string {
	var b strings.Builder
	b.Grow(len(text))

	for _, r := range text {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), unicode.IsSpace(r):
			b.WriteRune(unicode.ToLower(r)) // lowercase
		default:
			// skip punctuation/symbols
		}
	}
	return b.String()
}

// Tokenize splits text into tokens, removes empty tokens and stopwords, and applies stemming.
func Tokenize(text string, stopWords map[string]struct{}) []string {
	text = preprocessText(text)

	rawTokens := strings.Fields(text) // splits on any whitespace

	tokens := make([]string, 0, len(rawTokens))
	for _, tok := range rawTokens {
		if tok == "" {
			continue
		}

		// skip stopwords
		if _, found := stopWords[tok]; found {
			continue
		}

		// stem token
		stemmed := porterstemmer.StemString(tok)
		tokens = append(tokens, stemmed)
	}

	return tokens
}
