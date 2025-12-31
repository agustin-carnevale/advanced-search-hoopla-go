/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package hybrid

import (
	"fmt"
	"strconv"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

// normalizeCmd represents the normalize command
var normalizeCmd = &cobra.Command{
	Use:     "normalize <inputs...>",
	Short:   "Normalizes scores using the min-max normalization",
	Example: "normalize 1.2 3.2 6.4 3.0 0.3 2.0",
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

		results := methods.Normalize(inputs)

		for _, score := range results {
			fmt.Printf("* %.4f\n", score)
		}
	},
}

func init() {
	HybridCmd.AddCommand(normalizeCmd)
}
