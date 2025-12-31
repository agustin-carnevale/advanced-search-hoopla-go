/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newChunkCmd() *cobra.Command {
	var chunkSize int
	var overlap int

	cmd := &cobra.Command{
		Use:   "chunk <text> [--chunkSize <int>] [--overlap <int>]",
		Short: "Split long text into smaller pieces for embedding",
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

	cmd.Flags().IntVar(&chunkSize, "chunkSize", 200, "Specify the chunk size in words [default: 200]")
	cmd.Flags().IntVar(&overlap, "overlap", 0, "Specify number of words to overlap between chunks [default: 0]")

	return cmd
}

func init() {
	chunkCmd := newChunkCmd()
	SemanticCmd.AddCommand(chunkCmd)
}
