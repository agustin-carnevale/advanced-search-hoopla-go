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

// embedTextCmd represents the embedText command
var embedTextCmd = &cobra.Command{
	Use:   "embedText <text>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ Please provide a text to embed.")
			return
		}

		text := args[0]

		ss, err := utils.NewSemanticSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create semantic search client: %v\n", err)
		}

		embedding, err := ss.EmbedText(text)
		if err != nil {
			log.Fatalf("❌ Failed to create embedding: %v\n", err)
		}

		fmt.Println("Text:", text)
		fmt.Println("First 3 dimensions:", embedding[:3])
		fmt.Println("Dimensions:", len(embedding))

	},
}

func init() {
	SemanticCmd.AddCommand(embedTextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// embedTextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// embedTextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
