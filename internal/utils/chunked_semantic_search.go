package utils

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
)

type ChunkSimilarityScore struct {
	MovieIdx int
	ChunkIdx int
	Score    float64
}

type ChunkMetadataFile struct {
	Chunks      []ChunkMetadata `json:"chunks"`
	TotalChunks int             `json:"total_chunks"`
}

type Embedding = []float64

type ChunkMetadata struct {
	MovieIdx    int `json:"movie_idx"`
	ChunkIdx    int `json:"chunk_idx"`
	TotalChunks int `json:"total_chunks"`
}

type ChunkedSemanticSearch struct {
	*SemanticSearch
	ChunksEmbeddings []Embedding
	ChunksMetadata   []ChunkMetadata
}

func NewChunkedSemanticSearch(modelName string) (*ChunkedSemanticSearch, error) {
	ss, err := NewSemanticSearch(modelName)
	if err != nil {
		return nil, err
	}

	return &ChunkedSemanticSearch{
		SemanticSearch:   ss,
		ChunksEmbeddings: make([]Embedding, 0),
		ChunksMetadata:   make([]ChunkMetadata, 0),
	}, nil
}

func (css *ChunkedSemanticSearch) BuildChunksEmbeddings() ([]Embedding, error) {
	chunks := make([]string, 0)
	for docIndex, doc := range css.Documents {
		if len(doc.Description) > 0 {
			descChunks := SemanticChunk(doc.Description, 4, 1)
			chunks = append(chunks, descChunks...)
			for chunkIndex := range descChunks {
				css.ChunksMetadata = append(css.ChunksMetadata, ChunkMetadata{
					MovieIdx:    docIndex,
					ChunkIdx:    chunkIndex,
					TotalChunks: len(descChunks),
				})
			}
		}
	}

	fmt.Println("Total chunks:", len(chunks))

	// Create embeddings for all chunks
	embeddings, err := css.createEmbeddingsParallel(chunks)
	if err != nil {
		return nil, err
	}
	css.ChunksEmbeddings = embeddings

	if err = css.saveChunksEmbeddings(); err != nil {
		return nil, fmt.Errorf("Error saving ChunksEmbeddings to file: %w", err)
	}

	if err = css.saveChunksMetadata(len(chunks)); err != nil {
		return nil, fmt.Errorf("Error saving ChunksMetadata to file: %w", err)
	}

	return css.ChunksEmbeddings, nil
}

func (css *ChunkedSemanticSearch) saveChunksMetadata(totalChunks int) error {
	file, err := os.Create(fs.ChunksMetadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // same as indent=2 in Python

	metadata := ChunkMetadataFile{
		Chunks:      css.ChunksMetadata,
		TotalChunks: totalChunks,
	}

	return encoder.Encode(metadata)
}

func (css *ChunkedSemanticSearch) saveChunksEmbeddings() error {
	file, err := os.Create(fs.ChunksEmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(css.ChunksEmbeddings)
}

func (css *ChunkedSemanticSearch) loadChunksEmbeddings() error {
	file, err := os.Open(fs.ChunksEmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&css.ChunksEmbeddings)
}

func (css *ChunkedSemanticSearch) loadChunksMetadata() error {
	file, err := os.Open(fs.ChunksMetadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data ChunkMetadataFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	css.ChunksMetadata = data.Chunks
	return nil
}

func (css *ChunkedSemanticSearch) LoadOrCreateChunksEmbeddings(docs []model.Movie) ([]Embedding, error) {
	css.Documents = docs
	for _, doc := range docs {
		css.DocumentMap[doc.ID] = doc
	}

	// Check if both files exist
	_, embErr := os.Stat(fs.ChunksEmbeddingsPath)
	_, metaErr := os.Stat(fs.ChunksMetadataPath)

	if embErr == nil && metaErr == nil {
		// Try loading embeddings
		if err := css.loadChunksEmbeddings(); err == nil {
			// Try loading metadata
			if err := css.loadChunksMetadata(); err == nil {
				fmt.Println("Loaded existing chunk embeddings + metadata from disk.")
				return css.ChunksEmbeddings, nil
			}
		}

		// If loading fails → rebuild both
		fmt.Println("⚠️ Existing files were corrupt. Rebuilding all embeddings...")
	} else {
		fmt.Println("No existing chunks. Building embeddings from scratch...")
	}

	// Create from scratch
	return css.BuildChunksEmbeddings()
}

func (css *ChunkedSemanticSearch) SearchChunked(query string, limit int) ([]SemanticSearchResult, error) {
	// todo: check chunks_embeddings are valid/correctly loaded

	queryEmbedding, err := css.EmbedText(query)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create embedding of the query: %v\n", err)
	}

	scores := make([]ChunkSimilarityScore, 0, len(css.ChunksEmbeddings))
	for i, chunkEmbedding := range css.ChunksEmbeddings {
		similarityScore := CosineSimilarity(queryEmbedding, chunkEmbedding)
		chunkMetadata := css.ChunksMetadata[i]
		scores = append(scores, ChunkSimilarityScore{
			MovieIdx: chunkMetadata.MovieIdx,
			ChunkIdx: chunkMetadata.ChunkIdx,
			Score:    similarityScore,
		})
	}

	moviesScoreMap := make(map[int]float64)
	for _, scoreItem := range scores {
		movieIdx := scoreItem.MovieIdx
		chunkScore := scoreItem.Score
		_, ok := moviesScoreMap[movieIdx]
		if ok {
			if chunkScore > moviesScoreMap[movieIdx] {
				moviesScoreMap[movieIdx] = chunkScore
			}
		} else {
			moviesScoreMap[movieIdx] = chunkScore
		}
	}

	movieScores := make([]SimilarityScore, 0, len(moviesScoreMap))
	for movieIdx, movieScore := range moviesScoreMap {
		movieScores = append(movieScores, SimilarityScore{
			Score: movieScore,
			Movie: css.Documents[movieIdx],
		})
	}

	sort.Slice(movieScores, func(i, j int) bool {
		return movieScores[i].Score > movieScores[j].Score
	})

	results := make([]SemanticSearchResult, 0, limit)
	for _, item := range movieScores[:limit] {
		results = append(results, SemanticSearchResult{
			DocID:       item.Movie.ID,
			Score:       item.Score,
			Title:       item.Movie.Title,
			Description: item.Movie.Description,
		})
	}

	return results, nil
}
