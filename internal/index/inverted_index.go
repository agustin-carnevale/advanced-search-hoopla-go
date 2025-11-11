package index

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"sort"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/tokenizer"
)

type SearchResult struct {
	DocID int
	Score float64
	Movie model.Movie // optional: include your movie struct here if needed
}

type InvertedIndex struct {
	Index           map[string]map[int]struct{} // term -> set of doc IDs
	DocMap          map[int]model.Movie         // docID -> movie
	TermFrequencies map[int]map[string]int      // docID -> term -> count
	DocLengths      map[int]int                 // docID -> docLength
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{
		Index:           make(map[string]map[int]struct{}),
		DocMap:          make(map[int]model.Movie),
		TermFrequencies: make(map[int]map[string]int),
		DocLengths:      make(map[int]int),
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
	idx.DocLengths[docID] = len(tokens)
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

func (idx *InvertedIndex) getAvgDocLength() float64 {
	count := len(idx.DocLengths)
	if count == 0 {
		return 0.0
	}

	sum := 0
	for _, value := range idx.DocLengths {
		sum += value
	}

	return float64(sum) / float64(count)

}

func (idx *InvertedIndex) GetTF(docID int, term string) int {
	stopWords, err := fs.LoadStopWords()
	if err != nil {
		log.Fatalf("Error loading stop words: could not tokenize term.")
		return 0
	}
	tokens := tokenizer.Tokenize(term, stopWords)

	if len(tokens) == 0 {
		return 0
	}

	if len(tokens) > 1 {
		log.Fatalf("Error at get_tf(): term has too many tokens.")
		return 0
	}

	t := tokens[0]
	tfMap, exists := idx.TermFrequencies[docID]
	if !exists {
		return 0
	}

	return tfMap[t]
}

func (idx *InvertedIndex) GetIDF(term string) float64 {
	docCount := len(idx.DocMap)

	stopWords, err := fs.LoadStopWords()
	if err != nil {
		log.Fatalf("Error loading stop words: could not tokenize term.")
		return 0
	}
	tokens := tokenizer.Tokenize(term, stopWords)

	if len(tokens) == 0 {
		return 0
	}

	if len(tokens) > 1 {
		log.Fatalf("Error at get_idf(): term has too many tokens.")
		return 0
	}

	t := tokens[0]

	termDocCount := len(idx.Index[t])

	return math.Log(float64(docCount+1) / float64(termDocCount+1))
}

func (idx *InvertedIndex) GetBM25IDF(term string) float64 {
	N := len(idx.DocMap)

	stopWords, err := fs.LoadStopWords()
	if err != nil {
		log.Fatalf("Failed to load stop words")
		return 0.0
	}

	tokens := tokenizer.Tokenize(term, stopWords)
	if len(tokens) != 1 {
		log.Fatalf("Error at get_bm25_idf(): term has not a single token")
		return 0.0
	}

	t := tokens[0]
	docIDs := idx.Index[t]
	df := len(docIDs)

	return math.Log((float64(N)-float64(df)+0.5)/(float64(df)+0.5) + 1)
}

func (idx *InvertedIndex) GetBM25TF(docID int, term string, k1 float64, b float64) float64 {
	tf := float64(idx.GetTF(docID, term))

	avgDocLength := idx.getAvgDocLength()
	docLength := idx.DocLengths[docID]

	// Length normalization factor
	lengthNorm := 1 - b + b*(float64(docLength)/avgDocLength)

	// Apply to term frequency
	bm25tf := (tf * (k1 + 1)) / (tf + k1*lengthNorm)

	return bm25tf
}

func (idx *InvertedIndex) bm25Score(docID int, term string) float64 {
	bm25idf := idx.GetBM25IDF(term)
	bm25tf := idx.GetBM25TF(docID, term, 1.5, 0.75)

	return bm25idf * bm25tf
}

func (idx *InvertedIndex) Bm25Search(query string, limit int) []SearchResult {
	stopWords, err := fs.LoadStopWords()
	if err != nil {
		log.Fatalf("Error loading stop words: could not tokenize term.")
	}
	qTokens := tokenizer.Tokenize(query, stopWords)

	// Preallocate results slice directly
	results := make([]SearchResult, 0, len(idx.DocMap))
	for docID, movie := range idx.DocMap {
		var score float64
		for _, t := range qTokens {
			score += idx.bm25Score(docID, t)
		}

		results = append(results, SearchResult{
			DocID: docID,
			Score: score,
			Movie: movie,
		})
	}

	// Sort by Score DESC
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit < len(results) {
		results = results[:limit]
	}

	return results
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

func (idx *InvertedIndex) Load() error {
	f, err := os.Open(fs.IndexPath)
	if err != nil {
		return fmt.Errorf("failed to open index file: %w", err)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)

	if err := decoder.Decode(idx); err != nil {
		return fmt.Errorf("failed to decode index: %w", err)
	}

	return nil
}
