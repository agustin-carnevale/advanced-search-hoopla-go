/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

func newSemanticChunkCmd() *cobra.Command {
	var maxChunkSize int
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

			chunks := utils.SemanticChunk(text, maxChunkSize, overlap)

			for i, chunk := range chunks {
				fmt.Printf("%d. %s\n", i+1, chunk)
			}
		},
	}

	cmd.Flags().IntVar(&maxChunkSize, "maxChunkSize", 4, "Specify the chunk size in sentences")
	cmd.Flags().IntVar(&overlap, "overlap", 0, "Specify number of sentences to overlap between chunks")

	return cmd
}

func init() {
	semanticChunkCmd := newSemanticChunkCmd()
	SemanticCmd.AddCommand(semanticChunkCmd)
}
