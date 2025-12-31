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

// bm25idfCmd represents the bm25idf command
var bm25idfCmd = &cobra.Command{
	Use:   "bm25idf <term>",
	Short: "Get BM25 IDF score for a given term",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please provide a term.")
			return
		}
		term := args[0]

		// load index
		idx := index.NewInvertedIndex()
		if err := idx.Load(); err != nil {
			log.Fatalf("❌ Failed to load index: %v\n", err)
		}

		bm25idf := idx.GetBM25IDF(term)

		fmt.Printf("BM25 IDF score of '%s': %.2f\n", term, bm25idf)
	},
}

func init() {
	KeywordCmd.AddCommand(bm25idfCmd)
}
