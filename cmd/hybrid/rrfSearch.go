/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package hybrid

import (
	"context"
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

func newRRFSearchCmd() *cobra.Command {
	var limit int
	var k int
	var enhance string

	cmd := &cobra.Command{
		Use:   "rrfSearch <query> [--limit <int>] [--k <int>] [--enhance <spell|rewrite|expand>]",
		Short: "Reciprocal Rank Fusion search combining both keyword and semantic.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if enhance == "" {
				return nil
			}

			switch enhance {
			case "spell", "rewrite", "expand":
				return nil
			default:
				return fmt.Errorf(
					"invalid value for --enhance: %q (allowed: spell, rewrite, expand)",
					enhance,
				)
			}
		},
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

			if enhance != "" {
				ctx := context.Background()
				enhancedQuery, err := llms.PreProcessQuery(ctx, query, enhance)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				fmt.Printf("Enhanced query (%s): '%s' -> '%s'\n", enhance, query, enhancedQuery)
				query = enhancedQuery
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
	cmd.Flags().StringVar(&enhance, "enhance", "", "Query enhancement method. [choices: spell|rewrite|expand]")

	return cmd

}

func init() {
	rrfSearchCmd := newRRFSearchCmd()
	rrfSearchCmd.RegisterFlagCompletionFunc(
		"enhance",
		cobra.FixedCompletions(
			[]string{"spell", "rewrite", "expand"},
			cobra.ShellCompDirectiveNoFileComp,
		),
	)
	HybridCmd.AddCommand(rrfSearchCmd)
}
