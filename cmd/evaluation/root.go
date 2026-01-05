package evaluation

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

var limit int

var EvaluationCmd = &cobra.Command{
	Use:     "evaluation [--limit <int>]",
	Aliases: []string{"eval"},
	Short:   "Evaluation of the golden dataset",
	Run: func(cmd *cobra.Command, args []string) {
		testCases, err := fs.LoadGoldenDataset()
		if err != nil {
			log.Fatalf("❌ Failed to load golden dataset: %v\n", err)
		}

		hs, err := methods.NewHybridSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create hybrid search client: %v\n", err)
		}

		fmt.Printf("k=%d\n\n", limit)
		for i, testCase := range testCases {
			query := testCase.Query
			k := 60
			rrfSearchResults, err := hs.RRFSearch(query, k, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perform rrf search: %v\n", err)
			}

			retrievedTitles := make([]string, len(rrfSearchResults))
			for i, r := range rrfSearchResults {
				retrievedTitles[i] = r.Title
			}

			relevantCount := 0
			for _, title := range testCase.RelevantDocs {
				if slices.Contains(retrievedTitles, title) {
					relevantCount++
				}
			}

			precision := float64(relevantCount) / float64(len(retrievedTitles))

			fmt.Printf("Test Case %d\n", i+1)
			fmt.Printf("- Query: %s\n", query)
			fmt.Printf("\t- Precision@%d: %.4f\n", limit, precision)
			fmt.Printf("\t- Retrieved: %s\n", strings.Join(retrievedTitles, ", "))
			fmt.Printf("\t- Relevant: %s\n\n", strings.Join(testCase.RelevantDocs, ", "))
		}

	},
}

func init() {
	EvaluationCmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
}
