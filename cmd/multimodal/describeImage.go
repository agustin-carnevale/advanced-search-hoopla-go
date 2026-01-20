package multimodal

import (
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/spf13/cobra"
)

func newDescribeImageCmd() *cobra.Command {
	var imagePath string

	cmd := &cobra.Command{
		Use:   "describeImage <query> --imagePath <string>",
		Short: "Use llm to describe image from binary",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				fmt.Println("❌ Please provide a query.")
				return
			}
			query := args[0]

			// Image file into bytes
			img, err := os.ReadFile(imagePath)
			if err != nil {
				log.Fatalf("❌ Failed to read image file: %v\n", err)
			}

			// MIME type from file extension
			ext := filepath.Ext(imagePath)
			mime := mime.TypeByExtension(ext)
			if mime == "" {
				mime = "image/jpeg"
			}

			newQuery, totalTokens, err := llms.RewriteQueryFromImage(context.Background(), query, img, mime)
			if err != nil {
				log.Fatalf("❌ Failed calling llm with multimodal prompt: %v\n", err)
			}

			fmt.Printf("Rewritten query: %s\n", newQuery)
			fmt.Printf("Total tokens:\t%d\n", totalTokens)

		}}

	cmd.Flags().StringVar(&imagePath, "imagePath", "", "Path to the image.")

	return cmd
}

func init() {
	describeImageCmd := newDescribeImageCmd()
	describeImageCmd.MarkFlagRequired("imagePath")
	MultimodalCmd.AddCommand(describeImageCmd)
}
