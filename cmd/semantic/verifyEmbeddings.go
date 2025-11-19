/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

// verifyEmbeddingsCmd represents the verifyEmbeddings command
var verifyEmbeddingsCmd = &cobra.Command{
	Use:   "verifyEmbeddings",
	Short: "Verifies embeddings exist if not creates them",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := utils.NewSemanticSearch("nomic-embed-text")
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
