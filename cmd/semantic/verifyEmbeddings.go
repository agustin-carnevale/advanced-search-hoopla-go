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

var verifyEmbeddingsCmd = &cobra.Command{
	Use:   "verifyEmbeddings",
	Short: "Verifies embeddings exist if not creates them",
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := methods.NewSemanticSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create semantic search client: %v\n", err)
		}

		moviesDocs, err := fs.LoadMovies()
		if err != nil {
			log.Fatalf("❌ Failed to load movies: %v\n", err)
		}

		embeddings, err := ss.LoadOrCreateEmbeddings(moviesDocs)
		if err != nil {
			log.Fatalf("❌ Failed to load or generate embeddings: %v\n", err)
		}

		fmt.Printf("Number of docs: %d\n", len(moviesDocs))
		if len(embeddings) > 0 {
			fmt.Printf("Embeddings shape: %d vectors in %d dimensions\n", len(embeddings), len(embeddings[0]))
		}

	},
}

func init() {
	SemanticCmd.AddCommand(verifyEmbeddingsCmd)
}
