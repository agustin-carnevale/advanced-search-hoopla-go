/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package keyword

import (
	"fmt"
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔨 Starting build...")

		idx := index.NewInvertedIndex()

		if err := idx.Build(); err != nil {
			log.Fatalf("❌ Failed to build index: %v\n", err)
		}

		if err := idx.Save(); err != nil {
			log.Fatalf("❌ Failed to save index: %v\n", err)
		}

		fmt.Println("✅ Index build completed and saved!")
	},
}

func init() {
	KeywordCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
