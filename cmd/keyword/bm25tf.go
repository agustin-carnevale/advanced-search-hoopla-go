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

// bm25tfCmd represents the bm25tf command
var bm25tfCmd = &cobra.Command{
	Use:   "bm25tf <docID> <term> [k1] [b]",
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

		k1 := 1.5
		if len(args) > 2 {
			k1String := args[2]
			k1, err = strconv.ParseFloat(k1String, 64)
			if err != nil {
				log.Fatalf("❌ k1 should be a float: %v\n", err)
			}
		}

		b := 0.75
		if len(args) > 3 {
			bString := args[3]
			b, err = strconv.ParseFloat(bString, 64)
			if err != nil {
				log.Fatalf("❌ b should be a float: %v\n", err)
			}
		}

		// load index
		idx := index.NewInvertedIndex()
		if err := idx.Load(); err != nil {
			log.Fatalf("❌ Failed to load index: %v\n", err)
		}

		bm25tf := idx.GetBM25TF(docID, term, k1, b)

		fmt.Printf("BM25 TF score of '%s' in document '%d': %.2f\n", term, docID, bm25tf)

	},
}

func init() {
	KeywordCmd.AddCommand(bm25tfCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bm25tfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bm25tfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
