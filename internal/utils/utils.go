package utils

import (
	"fmt"
	"strings"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
)

func ResultsListToStr(results []methods.RRFSearchResult) string {
	parts := make([]string, len(results))
	for i, doc := range results {
		parts[i] = fmt.Sprintf(
			"ID: %d\nTitle: %s\nDescription: %s",
			doc.DocID,
			doc.Title,
			doc.Description,
		)
	}
	resultListStr := strings.Join(parts, "\n\n")

	return resultListStr
}
