/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package hybrid

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

func newWeightedSearchCmd() *cobra.Command {
	var limit int
	var alpha float64

	cmd := &cobra.Command{
		Use:   "weightedSearch <query> [--limit <int>] [--alpha <float>]",
		Short: "Weighted search combining both keyword and semantic",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("❌ Please provide a query.")
				return
			}
			query := args[0]

			hs, err := methods.NewHybridSearch("nomic-embed-text")
			if err != nil {
				log.Fatalf("❌ Failed to create hybrid search client: %v\n", err)
			}

			results, err := hs.WeightedSearch(query, alpha, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perform weighted search: %v\n", err)
			}

			for i, result := range results {
				fmt.Printf("%d. %s\n", i+1, result.Title)
				fmt.Printf("\tHybrid Score: %.3f\n", result.HybridScore)
				fmt.Printf("\tBM25: %.3f, Semantic: %.3f\n", result.KeywordScore, result.SemanticScore)
				fmt.Printf("\t%s...\n\n", result.Description[:100])
			}
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
	cmd.Flags().Float64Var(&alpha, "alpha", 0.5, "Dynamically control the weighting between the two scores")

	return cmd

}

func init() {
	weightedSearchCmd := newWeightedSearchCmd()
	HybridCmd.AddCommand(weightedSearchCmd)
}
