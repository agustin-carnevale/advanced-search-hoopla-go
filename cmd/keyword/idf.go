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

var idfCmd = &cobra.Command{
	Use:   "idf <term>",
	Short: "Calculate the inverse document frequency.",
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

		idf := idx.GetIDF(term)

		fmt.Printf("Inverse document frequency of '%s': %.2f\n", term, idf)
	},
}

func init() {
	KeywordCmd.AddCommand(idfCmd)
}
