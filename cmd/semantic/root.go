package semantic

import "github.com/spf13/cobra"

var SemanticCmd = &cobra.Command{
	Use:     "semantic",
	Aliases: []string{"sm", "st"},
	Short:   "Semantic search",
}
