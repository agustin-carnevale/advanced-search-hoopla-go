package multimodal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/methods"
	"github.com/spf13/cobra"
)

func newImageSearchCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "imageSearch <imagePath>",
		Short: "Use RAG to summarize the results with an LLM.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				fmt.Println("❌ Please provide a query.")
				return
			}
			imagePath := args[0]

			// Image file into bytes
			_, err := os.ReadFile(imagePath)
			if err != nil {
				log.Fatalf("❌ Failed to read image file: %v\n", err)
			}

			mms, err := methods.NewMultimodalSearch()
			if err != nil {
				log.Fatalf("❌ Failed to create MultimodalSearch client: %v\n", err)
			}

			_, err = mms.LoadOrCreateEmbeddings()
			if err != nil {
				log.Fatalf("❌ Failed to load or generate embeddings: %v\n", err)
			}

			// fmt.Println("embeddings len:", len(embeddings))
			// if len(embeddings) > 0 {
			// 	fmt.Println("dimensions:", len(embeddings[0]))
			// }

			ctx := context.Background()
			limit := 5
			results, err := mms.ImageSearch(ctx, imagePath, limit)
			if err != nil {
				log.Fatalf("❌ Failed to perfom image search: %v\n", err)
			}

			fmt.Printf("Image search results for: %s\n", imagePath)
			for i, res := range results {
				fmt.Printf("%d. %s (similarity: %.3f)\n", i+1, res.Title, res.Score)
				// fmt.Printf("\t%s\n\n", res.Description)
			}
		}}

	return cmd
}

func init() {
	imageSearchCmd := newImageSearchCmd()
	MultimodalCmd.AddCommand(imageSearchCmd)
}
