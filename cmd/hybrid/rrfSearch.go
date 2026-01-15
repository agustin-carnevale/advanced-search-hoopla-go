/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package hybrid

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/cli"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/logging"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func newRRFSearchCmd() *cobra.Command {
	var limit int
	var k int
	var enhance string
	var rerankMethod string
	var debug bool
	var evaluate bool

	cmd := &cobra.Command{
		Use:   "rrfSearch <query> [--limit <int>] [--k <int>] [--enhance <spell|rewrite|expand>] [--rerankMethod <individual|batch|crossEncoder>]",
		Short: "Reciprocal Rank Fusion search combining both keyword and semantic.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.ValidateFlagEnum(enhance, "enhance", "spell", "rewrite", "expand"); err != nil {
				return err
			}
			if err := cli.ValidateFlagEnum(
				rerankMethod,
				"rerankMethod",
				"individual",
				"batch",
				"crossEncoder",
			); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				fmt.Println("❌ Please provide a query.")
				return
			}
			query := args[0]

			// Initialize logger
			logger := logging.New(debug)
			execCtx := logging.ExecutionContext{
				RunID:   uuid.New().String(),
				QueryID: uuid.New().String(),
			}

			logging.LogOriginalQuery(logger, execCtx, query)

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

				logging.LogEnhancedQuery(logger, execCtx, logging.EnhancedQueryLog{
					EnhancementType: enhance,
					OriginalQuery:   query,
					EnhancedQuery:   enhancedQuery,
				})

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

			// Log RRF candidates
			if debug {
				rrfLogs := make([]logging.RRFCandidateLog, len(results))
				for i, r := range results {
					rrfLogs[i] = logging.RRFCandidateLog{
						DocID:        strconv.Itoa(r.DocID),
						RRFScore:     r.RRFScore,
						BM25Rank:     r.KeywordRank,
						SemanticRank: r.SemanticRank,
					}
				}
				logging.LogRRFResults(logger, execCtx, rrfLogs)
			}

			// normalize into a single result type
			var finalResults []methods.RRFSearchReRankedResult
			if rerankMethod != "" {
				// apply LLM reranking
				fmt.Printf("Reranking top %d results using %s method...\n", len(results), rerankMethod)
				finalResults, err = methods.ReRankResults(query, results, rerankMethod)
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

			// Log Final Results
			if debug {
				finalLogs := make([]logging.FinalResultLog, len(finalResults))
				for i, r := range finalResults {
					finalLogs[i] = logging.FinalResultLog{
						DocID:      strconv.Itoa(r.DocID),
						FinalScore: r.ReRankScore,
						Position:   i + 1,
					}
				}
				logging.LogFinalResults(logger, execCtx, finalLogs)
			}

			// print top results
			printRRFResults(finalResults, limit, rerankMethod, query, k)

			// perform LLM evaluation of results
			if evaluate {
				ctx := context.Background()
				// create list of results as string
				parts := make([]string, limit)
				for i, doc := range finalResults[:limit] {
					parts[i] = fmt.Sprintf(
						"ID: %d\nTitle: %s\nDescription: %s",
						doc.DocID,
						doc.Title,
						doc.Description,
					)
				}
				resultListStr := strings.Join(parts, "\n\n")

				scores, err := llms.EvaluateResults(ctx, query, resultListStr)
				if err != nil {
					log.Fatalf("❌ Failed to perform results evaluation: %v\n", err)
				}
				printRRFEvaluationResults(finalResults, limit, scores)
			}

		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
	cmd.Flags().IntVar(&k, "k", 60, "Controls how much more weight we give to higher-ranked results vs lower-ranked ones.")
	cmd.Flags().StringVar(&enhance, "enhance", "", "Query enhancement method. [choices: spell|rewrite|expand]")
	cmd.Flags().StringVar(&rerankMethod, "rerankMethod", "", "Re-ranking method. [choices: individual|batch|crossEncoder]")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug logging")
	cmd.Flags().BoolVar(&evaluate, "evaluate", false, "Add LLM evaluation to the results")

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
	rrfSearchCmd.RegisterFlagCompletionFunc(
		"rerankMethod",
		cobra.FixedCompletions(
			[]string{"individual", "batch", "crossEncoder"},
			cobra.ShellCompDirectiveNoFileComp,
		),
	)
	HybridCmd.AddCommand(rrfSearchCmd)
}

func printRRFResults(
	results []methods.RRFSearchReRankedResult,
	limit int,
	rerankMethod string,
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
		if rerankMethod == "individual" {
			fmt.Printf("\tReRank Score: %.3f/10\n", result.ReRankScore)
		}
		if rerankMethod == "batch" {
			fmt.Printf("\tReRank Rank: %d\n", i+1)
		}
		if rerankMethod == "crossEncoder" {
			fmt.Printf("\tCross Encoder Score: %.3f\n", result.ReRankScore)
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

func printRRFEvaluationResults(
	results []methods.RRFSearchReRankedResult,
	limit int,
	scores []int,
) {
	if len(scores) != limit {
		fmt.Println("LLM evaluation gave inconsistent results.")
		return
	}

	fmt.Println("LLM evaluation results:")

	for i, result := range results[:limit] {
		fmt.Printf("%d. %s: %d/3\n", i+1, result.Title, scores[i])
	}

}
