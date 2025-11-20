/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func newSemanticChunkCmd() *cobra.Command {
	var chunkSize int
	var overlap int

	cmd := &cobra.Command{
		Use:   "semanticChunk <text> [--maxChunkSize <value>] [--overlap <value>]",
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

			// TODO: this is removing the punctuation (add it back?)
			sentenceRegex := regexp.MustCompile(`([.!?])\s+`)
			sentences := sentenceRegex.Split(text, -1)
			if len(sentences) == 0 {
				return
			}

			var chunks []string
			start := 0
			end := min(chunkSize, len(sentences))
			c := strings.Join(sentences[start:end], " ")
			chunks = append(chunks, c)

			start = min(chunkSize-overlap, len(sentences))
			for start+overlap < len(sentences) {
				end = min(start+chunkSize, len(sentences))
				c = strings.Join(sentences[start:end], " ")
				chunks = append(chunks, c)
				start = end - overlap
			}

			for i, chunk := range chunks {
				fmt.Printf("%d. %s\n", i+1, chunk)
			}
		},
	}

	cmd.Flags().IntVar(&chunkSize, "maxChunkSize", 4, "Specify the chunk size in sentences")
	cmd.Flags().IntVar(&overlap, "overlap", 0, "Specify number of sentences to overlap between chunks")

	return cmd
}

func init() {
	semanticChunkCmd := newSemanticChunkCmd()
	SemanticCmd.AddCommand(semanticChunkCmd)
}
