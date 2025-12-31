/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that the semantic model is correctly loaded and working",
	Run: func(cmd *cobra.Command, args []string) {
		ss, err := utils.NewSemanticSearch("nomic-embed-text")
		if err != nil {
			log.Fatalf("❌ Failed to create semantic search client: %v\n", err)
		}

		if err := ss.VerifyModel(); err != nil {
			log.Fatalf("❌ Failed to verify model: %v\n", err)
		}
	},
}

func init() {
	SemanticCmd.AddCommand(verifyCmd)
}
