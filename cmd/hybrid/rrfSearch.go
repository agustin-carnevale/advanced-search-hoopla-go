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

func newRRFSearchCmd() *cobra.Command {
	var limit int
	var k int

	cmd := &cobra.Command{
		Use:   "rrfSearch <query> [--limit <int>] [--k <int>]",
		Short: "RRF search combining both keyword and semantic",
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

			results, err := hs.RRFSearch(query, k, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perform weighted search: %v\n", err)
			}

			for i, result := range results {
				fmt.Printf("%d. %s\n", i+1, result.Title)
				fmt.Printf("\tRRF Score: %.3f\n", result.RRFScore)
				fmt.Printf("\tBM25 Rank: %d, Semantic Rank: %d\n", result.KeywordRank, result.SemanticRank)
				fmt.Printf("\t%s...\n\n", result.Description[:100])
			}
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
	cmd.Flags().IntVar(&k, "k", 60, "Controls how much more weight we give to higher-ranked results vs lower-ranked ones.")

	return cmd

}

func init() {
	rrfSearchCmd := newRRFSearchCmd()
	HybridCmd.AddCommand(rrfSearchCmd)
}
