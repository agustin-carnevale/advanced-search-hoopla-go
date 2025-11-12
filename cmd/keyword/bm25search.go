/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/spf13/cobra"
)

// bm25searchCmd represents the bm25search command
var bm25searchCmd = &cobra.Command{
	Use:   "bm25search query <limit>",
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
		limit := 5
		if len(args) > 1 {
			limitString := args[1]
			l, err := strconv.Atoi(limitString)
			if err != nil {
				log.Fatalf("❌ limit should be a int: %v\n", err)
				return
			}
			limit = l
		}

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

func init() {
	KeywordCmd.AddCommand(bm25searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bm25searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bm25searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
