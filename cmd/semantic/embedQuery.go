/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

var embedQueryCmd = &cobra.Command{
	Use:   "embedQuery <query>",
	Short: "A brief description of your command",
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

		ss, err := utils.NewSemanticSearch("nomic-embed-text")
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
