/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/tokenizer"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
