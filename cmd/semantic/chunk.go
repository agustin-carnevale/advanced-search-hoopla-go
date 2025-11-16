/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package semantic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// chunkCmd represents the chunk command
var chunkCmd = &cobra.Command{
	Use:   "chunk <text> [chunkSize]",
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

		chunkSize := 200
		if len(args) > 1 {
			cs, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("❌ chunkSize should be an int: %v\n", err)
			}
			chunkSize = cs
		}

		words := strings.Fields(text)

		var chunks []string
		for i := 0; i < len(words); i += chunkSize {
			end := i + chunkSize
			if end > len(words) {
				end = len(words)
			}
			c := strings.Join(words[i:end], " ")
			chunks = append(chunks, c)
		}

		for i, chunk := range chunks {
			fmt.Printf("%d. %s\n", i+1, chunk)
		}
	},
}

func init() {
	SemanticCmd.AddCommand(chunkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chunkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chunkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
