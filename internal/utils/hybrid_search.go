package utils

import (
	"log"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
)

type HybridSearch struct {
	Idx *index.InvertedIndex
	Css *ChunkedSemanticSearch
}

func NewHybridSearch(modelName string) (*HybridSearch, error) {
	// Create and build an inverted index
	idx := index.NewInvertedIndex()
	if err := idx.Build(); err != nil {
		log.Fatalf("❌ Failed to build index: %v\n", err)
	}
	if err := idx.Save(); err != nil {
		log.Fatalf("❌ Failed to save index: %v\n", err)
	}

	// Create and build chunked semantic search
	css, err := NewChunkedSemanticSearch(modelName)
	if err != nil {
		return nil, err
	}
	moviesDocs, err := fs.LoadMovies()
	if err != nil {
		log.Fatalf("❌ Failed to load movies: %v\n", err)
	}
	_, err = css.LoadOrCreateChunksEmbeddings(moviesDocs)
	if err != nil {
		log.Fatalf("❌ Failed to load or generate embeddings: %v\n", err)
	}

	return &HybridSearch{
		Idx: idx,
		Css: css,
	}, nil
}

func (hs *HybridSearch) bm25Search(query string, limit int) ([]index.SearchResult, error) {
	err := hs.Idx.Load()
	if err != nil {
		return nil, err
	}
	results := hs.Idx.Bm25Search(query, limit)
	return results, nil
}

func (hs *HybridSearch) WeightedSearch(query string, alpha int, limit int) {
	log.Fatalf("❌ Weighted hybrid search is not implemented yet.")
}

func (hs *HybridSearch) RRFSearch(query string, k float64, limit int) {
	log.Fatalf("❌ RRF hybrid search is not implemented yet.")
}
