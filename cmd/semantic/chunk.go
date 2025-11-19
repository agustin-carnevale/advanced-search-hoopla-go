/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var chunkSize int
var overlap int

// chunkCmd represents the chunk command
var chunkCmd = &cobra.Command{
	Use:   "chunk <text> [--chunkSize] [--overlap]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("❌ You need to provide a text to chunk")
			return
		}
		text := args[0]
		fmt.Printf("Chunking %d characters\n", len(text))

		words := strings.Fields(text)
		if len(words) == 0 {
			return
		}

		var chunks []string
		start := 0
		end := min(chunkSize, len(words))
		c := strings.Join(words[start:end], " ")
		chunks = append(chunks, c)

		start = min(chunkSize-overlap, len(words))
		for start+overlap < len(words) {
			end = min(start+chunkSize, len(words))
			c = strings.Join(words[start:end], " ")
			chunks = append(chunks, c)
			start = end - overlap
		}

		for i, chunk := range chunks {
			fmt.Printf("%d. %s\n", i+1, chunk)
		}
	},
}

func init() {
	SemanticCmd.AddCommand(chunkCmd)

	chunkCmd.Flags().IntVar(&chunkSize, "chunk-size", 200, "Specify the chunk size in words")
	chunkCmd.Flags().IntVar(&overlap, "overlap", 0, "Specify number of words to overlap between chunks")
}
