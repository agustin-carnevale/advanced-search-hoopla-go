/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package hybrid

import (
	"context"
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

func newRRFSearchCmd() *cobra.Command {
	var limit int
	var k int
	var enhance string
	var rerankMethod string

	cmd := &cobra.Command{
		Use:   "rrfSearch <query> [--limit <int>] [--k <int>] [--enhance <spell|rewrite|expand>] [--rerankMethod]",
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

			hs, err := methods.NewHybridSearch("nomic-embed-text")
			if err != nil {
				log.Fatalf("❌ Failed to create hybrid search client: %v\n", err)
			}

			// Pre-process query
			if enhance != "" {
				ctx := context.Background()
				enhancedQuery, err := llms.PreProcessQuery(ctx, query, enhance)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				fmt.Printf("Enhanced query (%s): '%s' -> '%s'\n", enhance, query, enhancedQuery)
				query = enhancedQuery
			}

			// Set search limit
			searchLimit := limit
			if rerankMethod != "" {
				searchLimit = limit * 5
			}

			results, err := hs.RRFSearch(query, k, searchLimit)
			if err != nil {
				log.Fatalf("❌ Failed to perform rrf search: %v\n", err)
			}

			// normalize into a single result type
			var finalResults []methods.RRFSearchReRankedResult
			if rerankMethod != "" {
				// apply LLM reranking
				fmt.Printf("Reranking top %d results using %s method...\n", len(results), rerankMethod)
				finalResults, err = methods.ReRankResults(query, results)
				if err != nil {
					log.Fatalf("❌ Failed to perform re-ranking: %v\n", err)
				}
			} else {
				// wrap base results as RRFSearchReRankedResult (to keep it uniform)
				finalResults = make([]methods.RRFSearchReRankedResult, len(results))
				for i, r := range results {
					finalResults[i] = methods.RRFSearchReRankedResult{
						RRFSearchResult: r,
					}
				}
			}
			// print top results
			printRRFResults(finalResults, limit, rerankMethod != "", query, k)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
	cmd.Flags().IntVar(&k, "k", 60, "Controls how much more weight we give to higher-ranked results vs lower-ranked ones.")
	cmd.Flags().StringVar(&enhance, "enhance", "", "Query enhancement method. [choices: spell|rewrite|expand]")
	cmd.Flags().StringVar(&rerankMethod, "rerankMethod", "", "Re-ranking method. [choices: individual]")

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

func printRRFResults(
	results []methods.RRFSearchReRankedResult,
	limit int,
	showReRank bool,
	query string,
	k int,
) {
	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}
	if limit > len(results) {
		limit = len(results)
	}

	fmt.Printf("Reciprocal Rank Fusion Results for '%s' (k=%d):\n\n", query, k)

	for i, result := range results[:limit] {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		if showReRank {
			fmt.Printf("\tReRank Score: %.3f/10\n", result.ReRankScore)
		}
		fmt.Printf("\tRRF Score: %.3f\n", result.RRFScore)
		fmt.Printf("\tBM25 Rank: %d, Semantic Rank: %d\n", result.KeywordRank, result.SemanticRank)
		desc := result.Description
		if len(desc) > 100 {
			desc = desc[:100]
		}
		fmt.Printf("\t%s...\n\n", desc)
	}
}
