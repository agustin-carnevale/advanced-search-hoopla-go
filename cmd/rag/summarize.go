package rag

import (
	"context"
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

func newSummarizeCmd() *cobra.Command {
	var limit int
	var k int

	cmd := &cobra.Command{
		Use:   "summarize <query>",
		Short: "Use RAG to summarize the results with an LLM.",
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

			results, err := hs.RRFSearch(query, k, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perform rrf search: %v\n", err)
			}

			ctx := context.Background()
			resultsStr := utils.ResultsListToStr(results)

			llmSummary, err := llms.ResultsAugmentation(ctx, query, resultsStr)

			fmt.Println("Search Results:")
			for _, r := range results {
				fmt.Printf("\t- %s\n", r.Title)
			}

			fmt.Println("LLM Summary:")
			fmt.Println(llmSummary)

		}}

	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")
	cmd.Flags().IntVar(&k, "k", 60, "Controls how much more weight we give to higher-ranked results vs lower-ranked ones.")

	return cmd
}

func init() {
	summarizeCmd := newSummarizeCmd()
	RAGCmd.AddCommand(summarizeCmd)
}
