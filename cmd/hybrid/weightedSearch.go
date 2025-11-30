/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package hybrid

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

func newWeightedSearchCmd() *cobra.Command {
	var limit int
	var alpha float64

	cmd := &cobra.Command{
		Use:   "weightedSearch <query> [--limit <int>] [--alpha <float>]",
		Short: "Weighted search combining both keyword and semantic",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("❌ Please provide a query to embed.")
				return
			}
			query := args[0]

			hs, err := utils.NewHybridSearch("nomic-embed-text")
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
