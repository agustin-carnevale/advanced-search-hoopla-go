package index

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/tokenizer"
)

type InvertedIndex struct {
	Index           map[string]map[int]struct{} // term -> set of doc IDs
	DocMap          map[int]model.Movie         // docID -> movie
	TermFrequencies map[int]map[string]int      // docID -> term -> count
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index:           make(map[string]map[int]struct{}),
		DocMap:          make(map[int]model.Movie),
		TermFrequencies: make(map[int]map[string]int),
	}
}

func (idx *InvertedIndex) addDocument(docID int, text string, stopWords map[string]struct{}) {
	tokens := tokenizer.Tokenize(text, stopWords)

	tf := make(map[string]int)
	for _, t := range tokens {
		tf[t]++

		// if term doesn't exist, initialize a "set"
		if _, exists := idx.Index[t]; !exists {
			idx.Index[t] = make(map[int]struct{})
		}
		idx.Index[t][docID] = struct{}{}
	}
	idx.TermFrequencies[docID] = tf
}

func (idx *InvertedIndex) GetDocuments(term string) []model.Movie {
	docIDs, ok := idx.Index[term]
	if !ok {
		return []model.Movie{}
	}

	movies := make([]model.Movie, 0, len(docIDs))
	for id := range docIDs {
		if movie, ok := idx.DocMap[id]; ok {
			movies = append(movies, movie)
		}
	}

	// Sort by ID (optional)
	slices.SortFunc(movies, func(a, b model.Movie) int {
		return a.ID - b.ID
	})

	return movies
}

func (idx *InvertedIndex) GetTF(docID int, term string) int {
	tfMap, exists := idx.TermFrequencies[docID]
	if !exists {
		return 0
	}

	return tfMap[term]
}

func (idx *InvertedIndex) GetIDF(term string) float64 {
	docCount := len(idx.DocMap)
	termDocCount := len(idx.Index[term])

	return math.Log(float64(docCount+1) / float64(termDocCount+1))
}

func (idx *InvertedIndex) GetBM25IDF(term string) float64 {
	N := len(idx.DocMap)
	df := len(idx.Index[term])

	return math.Log((float64(N)-float64(df)+0.5)/(float64(df)+0.5) + 1)
}

func (idx *InvertedIndex) Build() error {
	movies, err := fs.LoadMovies()
	if err != nil {
		return err
	}

	// map[string]struct{}
	stopWords, err := fs.LoadStopWords()
	if err != nil {
		return err
	}

	for _, movie := range movies {
		idx.DocMap[movie.ID] = movie
		text := movie.Title + " " + movie.Description
		idx.addDocument(movie.ID, text, stopWords)
	}

	return nil
}

func (idx *InvertedIndex) Save() error {
	// Make sure the cache dir exists
	if err := os.MkdirAll(fs.CacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create cache dir: %w", err)
	}

	f, err := os.Create(fs.IndexPath)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(idx); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	return nil
}

func Load() (*InvertedIndex, error) {
	f, err := os.Open(fs.IndexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}
	defer f.Close()

	var idx InvertedIndex
	decoder := gob.NewDecoder(f)

	if err := decoder.Decode(&idx); err != nil {
		return nil, fmt.Errorf("failed to decode index: %w", err)
	}

	return &idx, nil
}
