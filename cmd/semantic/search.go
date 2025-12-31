/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query> [--limit <int>]",
		Short: "Semantic search for query among all documents/movies",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("❌ Please provide a query to embed.")
				return
			}
			query := args[0]

			ss, err := methods.NewSemanticSearch("nomic-embed-text")
			if err != nil {
				log.Fatalf("❌ Failed to create semantic search client: %v\n", err)
			}

			moviesDocs, err := fs.LoadMovies()
			if err != nil {
				log.Fatalf("❌ Failed to load movies: %v\n", err)
			}

			_, err = ss.LoadOrCreateEmbeddings(moviesDocs)
			if err != nil {
				log.Fatalf("❌ Failed to load or generate embeddings: %v\n", err)
			}

			results, err := ss.Search(query, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perform semantic search: %v\n", err)
			}

			for i, result := range results {
				fmt.Printf("%d. %s (score: %.4f)\n", i+1, result.Title, result.Score)
				fmt.Printf("   %s ...\n\n", result.Description[:100])
			}

		},
	}
	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results [default: 5]")

	return cmd

}

func init() {
	searchCmd := newSearchCmd()
	SemanticCmd.AddCommand(searchCmd)
}
