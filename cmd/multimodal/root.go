package multimodal

import "github.com/spf13/cobra"

var MultimodalCmd = &cobra.Command{
	Use:     "multimodal",
	Aliases: []string{"mm"},
	Short:   "Multimodal commands",
}
