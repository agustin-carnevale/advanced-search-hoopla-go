/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build project inverted index.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔨 Starting build...")

		idx := index.NewInvertedIndex()

		if err := idx.Build(); err != nil {
			log.Fatalf("❌ Failed to build index: %v\n", err)
		}

		if err := idx.Save(); err != nil {
			log.Fatalf("❌ Failed to save index: %v\n", err)
		}

		fmt.Println("✅ Index build completed and saved!")
	},
}

func init() {
	KeywordCmd.AddCommand(buildCmd)
}
