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

var embedChunksCmd = &cobra.Command{
	Use:   "embedChunks",
	Short: "Verifies chunked embeddings exist if not creates them",
	Run: func(cmd *cobra.Command, args []string) {
		css, err := methods.NewChunkedSemanticSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create chunked semantic search client: %v\n", err)
		}

		moviesDocs, err := fs.LoadMovies()
		if err != nil {
			log.Fatalf("❌ Failed to load movies: %v\n", err)
		}

		embeddings, err := css.LoadOrCreateChunksEmbeddings(moviesDocs)
		if err != nil {
			log.Fatalf("❌ Failed to load or generate embeddings: %v\n", err)
		}

		fmt.Printf("Generated %d chunked embeddings\n", len(embeddings))
	},
}

func init() {
	SemanticCmd.AddCommand(embedChunksCmd)

}
