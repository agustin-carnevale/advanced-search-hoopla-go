/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"
	"strconv"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/spf13/cobra"
)

// tfidfCmd represents the tfidf command
var tfidfCmd = &cobra.Command{
	Use:   "tfidf <docID> <term>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println("❌ Please provide a docID and a term.")
			return
		}
		docIDString := args[0]
		docID, err := strconv.Atoi(docIDString)
		if err != nil {
			log.Fatalf("❌ docID should be an int: %v\n", err)
		}

		term := args[1]

		// load index
		idx := index.NewInvertedIndex()
		if err := idx.Load(); err != nil {
			log.Fatalf("❌ Failed to load index: %v\n", err)
		}

		tf := idx.GetTF(docID, term)
		idf := idx.GetIDF(term)

		tf_idf := float64(tf) * idf

		fmt.Printf("TF-IDF score of '%s' in document '%d': %.2f\n", term, docID, tf_idf)
	},
}

func init() {
	KeywordCmd.AddCommand(tfidfCmd)
}
