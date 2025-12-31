/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/tokenizer"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Inverted index based basic search",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("❌ Please provide a search query.")
			return
		}

		query := args[0]
		// query := strings.Join(args, " ") // support multi-word queries

		idx := index.NewInvertedIndex()

		if err := idx.Load(); err != nil {
			log.Fatalf("❌ Failed to load index: %v\n", err)
		}

		stopWords, err := fs.LoadStopWords()
		if err != nil {
			log.Fatalf("❌ Failed to load stop words: %v\n", err)
		}

		queryTokens := tokenizer.Tokenize(query, stopWords)

		results := []model.Movie{}
		for _, t := range queryTokens {
			docs := idx.GetDocuments(t)
			results = append(results, docs...)
			if len(results) >= 5 {
				break
			}
		}

		if len(results) == 0 {
			fmt.Println("No results found.")
			return
		}

		if len(results) > 5 {
			results = results[:5]
		}

		for i, movie := range results {
			fmt.Printf("%d. %s\n", i+1, movie.Title)
		}

	},
}

func init() {
	KeywordCmd.AddCommand(searchCmd)
}
