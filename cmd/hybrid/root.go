package hybrid

import "github.com/spf13/cobra"

var HybridCmd = &cobra.Command{
	Use:     "hybrid",
	Aliases: []string{"hb"},
	Short:   "Hybrid search commands",
}
