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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bm25idfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bm25idfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
