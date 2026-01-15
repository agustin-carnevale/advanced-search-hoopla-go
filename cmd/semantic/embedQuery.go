/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

var embedQueryCmd = &cobra.Command{
	Use:   "embedQuery <query>",
	Short: "Create embedding for the given query",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ Please provide a query.")
			return
		}
		query := args[0]

		ss, err := methods.NewSemanticSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create semantic search client: %v\n", err)
		}

		embedding, err := ss.EmbedText(query)
		if err != nil {
			log.Fatalf("❌ Failed to create embedding of the query: %v\n", err)
		}

		fmt.Println("Query:", query)
		fmt.Println("First 5 dimensions:", embedding[:5])
		fmt.Println("Shape (dimensions):", len(embedding))
	},
}

func init() {
	SemanticCmd.AddCommand(embedQueryCmd)
}
