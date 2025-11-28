/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package hybrid

import (
	"fmt"
	"strconv"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/utils"
	"github.com/spf13/cobra"
)

// normalizeCmd represents the normalize command
var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		inputs := make([]float64, 0, len(args))
		for _, a := range args {
			val, err := strconv.ParseFloat(a, 64)
			if err != nil {
				fmt.Printf("❌ invalid float %q: %v\n", a, err)
				return
			}
			inputs = append(inputs, val)
		}

		results := utils.Normalize(inputs)

		for _, score := range results {
			fmt.Printf("* %.4f\n", score)
		}
	},
}

func init() {
	HybridCmd.AddCommand(normalizeCmd)
}
