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

func NewBm25TFCmd() *cobra.Command {
	// Define vars for flags
	var k1 float64
	var b float64

	cmd := &cobra.Command{
		Use:   "bm25tf <docID> <term> [--k1 <float>] [--b <float>]",
		Short: "Get BM25 TF score for a given document ID and term",
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

			bm25tf := idx.GetBM25TF(docID, term, k1, b)

			fmt.Printf("BM25 TF score of '%s' in document '%d': %.2f\n", term, docID, bm25tf)

		},
	}

	cmd.Flags().Float64Var(&k1, "k1", 1.5, "Define a k1 parameter [default: 1.5]")
	cmd.Flags().Float64Var(&b, "b", 0.75, "Define a b parameter [default: 0.75]")

	return cmd

}
func init() {
	bm25tfCmd := NewBm25TFCmd()
	KeywordCmd.AddCommand(bm25tfCmd)
}
