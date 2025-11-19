/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"
	"time"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/spf13/cobra"
)

func newBm25SearchCmd() *cobra.Command {
	// Define vars for flags
	var limit int

	cmd := &cobra.Command{
		Use:   "bm25search query [--limit]",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("❌ Please provide a query to search for.")
				return
			}
			query := args[0]

			// load index
			idx := index.NewInvertedIndex()
			if err := idx.Load(); err != nil {
				log.Fatalf("❌ Failed to load index: %v\n", err)
			}

			// Benchmark (measure time) just for testing purposes
			start := time.Now()

			results := idx.Bm25Search(query, limit)

			elapsed := time.Since(start)
			fmt.Printf("Bm25Search execution time: %s\n", elapsed)

			for i, doc := range results {
				fmt.Printf("%d. (%d) %s - Score: %.2f\n", i+1, doc.DocID, doc.Movie.Title, doc.Score)
			}

		},
	}

	cmd.Flags().IntVar(&limit, "limit", 5, "Limit the amount of results")

	return cmd
}

func init() {
	bm25searchCmd := newBm25SearchCmd()
	KeywordCmd.AddCommand(bm25searchCmd)
}
