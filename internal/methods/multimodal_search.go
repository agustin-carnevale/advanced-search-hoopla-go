package methods

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"mime"
	"os"
	"path/filepath"
	"sort"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/llms"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
)

type ImageSearchResult struct {
	ID          int
	Title       string
	Description string
	Score       float32
}

type similaryScore struct {
	idx   int
	score float32
}

type MultimodalSearch struct {
	Documents      []model.Movie
	DocsEmbeddings [][]float32
}

func NewMultimodalSearch() (*MultimodalSearch, error) {
	docs, err := fs.LoadMovies()
	if err != nil {
		log.Fatalf("❌ Failed to load movies: %v\n", err)
	}

	return &MultimodalSearch{
		Documents:      docs,
		DocsEmbeddings: make([][]float32, 0),
	}, nil
}

func (mms *MultimodalSearch) LoadOrCreateEmbeddings() ([][]float32, error) {
	// Check if file exist
	_, embErr := os.Stat(fs.MultimodalEmbeddingsPath)

	if embErr == nil {
		// Try loading embeddings
		if err := mms.loadDocsEmbeddings(); err == nil {
			// Check dimensions mismatch
			if len(mms.DocsEmbeddings) > 0 && len(mms.DocsEmbeddings[0]) != 768 {
				fmt.Printf("Dimension mismatch (got %d, expected 768). Rebuilding all embeddings...\n", len(mms.DocsEmbeddings[0]))
			} else {
				fmt.Println("Loaded existing chunk embeddings + metadata from disk.")
				return mms.DocsEmbeddings, nil
			}
		} else {
			// If loading fails → rebuild
			fmt.Println("Existing files were corrupt. Rebuilding all embeddings...")
		}
	} else {
		fmt.Println("No existing file. Building embeddings from scratch...")
	}

	// Create from scratch
	return mms.BuildEmbeddings()
}

func (mms *MultimodalSearch) loadDocsEmbeddings() error {
	file, err := os.Open(fs.MultimodalEmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&mms.DocsEmbeddings)
}

func (mms *MultimodalSearch) saveEmbeddings() error {
	file, err := os.Create(fs.MultimodalEmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(mms.DocsEmbeddings)
}

func (mms *MultimodalSearch) BuildEmbeddings() ([][]float32, error) {
	// Build strings: "title: description"
	docsAsStrings := make([]string, len(mms.Documents))
	for i, doc := range mms.Documents {
		// mms.DocumentMap[doc.ID] = doc
		docsAsStrings[i] = fmt.Sprintf("%s: %s", doc.Title, doc.Description)
	}

	// Create embeddings for all dcos
	embeddings, err := llms.BuildMultimodalEmbeddingsBatch(context.Background(), mms.Documents)
	if err != nil {
		return nil, err
	}

	mms.DocsEmbeddings = embeddings

	if err = mms.saveEmbeddings(); err != nil {
		return nil, fmt.Errorf("Error saving ChunksEmbeddings to file: %w", err)
	}

	return mms.DocsEmbeddings, nil
}

func (mms *MultimodalSearch) ImageSearch(ctx context.Context, imagePath string, limit int) ([]ImageSearchResult, error) {

	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	// MIME type from file extension
	ext := filepath.Ext(imagePath)
	mime := mime.TypeByExtension(ext)
	if mime == "" {
		mime = "image/jpeg"
	}

	imageEmbedding, err := llms.EmbedImage(ctx, imageBytes, mime)
	if err != nil {
		return nil, fmt.Errorf("Error creating image embedding: %v", err)
	}

	var scores []similaryScore
	for i, emb := range mms.DocsEmbeddings {
		score := CosineSimilarityFloat32(imageEmbedding, emb)
		scores = append(scores, similaryScore{i, score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	var results []ImageSearchResult
	for _, s := range scores[:limit] {
		m := mms.Documents[s.idx]
		results = append(results, ImageSearchResult{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Score:       s.score,
		})
	}

	return results, nil
}

func CosineSimilarityFloat32(a, b []float32) float32 {
	var dot, na, nb float32
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / float32(math.Sqrt(float64(na*nb)))
}
