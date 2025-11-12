package keyword

import "github.com/spf13/cobra"

var KeywordCmd = &cobra.Command{
	Use:     "keyword",
	Aliases: []string{"kw", "kv"},
	Short:   "Keyword-based search commands",
}
