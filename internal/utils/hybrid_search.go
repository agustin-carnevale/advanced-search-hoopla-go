package utils

import (
	"fmt"
	"log"
	"slices"
	"sort"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/index"
)

type ScoredDoc struct {
	DocID  int
	Scores *CombinedScores
}
type CombinedScores struct {
	HybridScore   float64
	KeywordScore  float64
	SemanticScore float64
}

type WeightedSearchResult struct {
	DocID         int
	Title         string
	Description   string
	HybridScore   float64
	KeywordScore  float64
	SemanticScore float64
}
type RankedDoc struct {
	DocID int
	Ranks *CombinedRanks
}
type CombinedRanks struct {
	RRFScore     float64
	KeywordRank  int
	SemanticRank int
}

type RRFSearchResult struct {
	DocID        int
	Title        string
	Description  string
	RRFScore     float64
	KeywordRank  int
	SemanticRank int
}

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

func (hs *HybridSearch) WeightedSearch(query string, alpha float64, limit int) ([]WeightedSearchResult, error) {
	searchLimit := min(limit*500, len(hs.Css.Documents))

	// keyword search
	keywordResults, err := hs.bm25Search(query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform bm25Search: %v\n", err)
	}

	scores := make([]float64, len(keywordResults))
	for i, v := range keywordResults {
		scores[i] = v.Score
	}
	normalizedKeywordScores := Normalize(scores)

	// semantic search
	semanticResults, err := hs.Css.SearchChunked(query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform SearchChunked: %v\n", err)
	}

	scores = make([]float64, len(semanticResults))
	for i, v := range semanticResults {
		scores[i] = v.Score
	}
	normalizedSemanticScores := Normalize(scores)

	// Collect all normalized scores by DocID
	combinedScores := make(map[int]*CombinedScores, len(keywordResults)+len(semanticResults))

	// Fill keyword scores
	for i, r := range keywordResults {
		combinedScores[r.DocID] = &CombinedScores{
			KeywordScore: normalizedKeywordScores[i],
		}
	}

	// Fill semantic scores
	for i, r := range semanticResults {
		cs, ok := combinedScores[r.DocID]
		if !ok {
			cs = &CombinedScores{}
			combinedScores[r.DocID] = cs
		}
		cs.SemanticScore = normalizedSemanticScores[i]
	}

	// Compute hybrid scores
	for _, cs := range combinedScores {
		cs.HybridScore = HybridScore(cs.KeywordScore, cs.SemanticScore, alpha)
	}

	// map -> slice for sorting
	scoredDocs := make([]ScoredDoc, 0, len(combinedScores))
	for docID, scores := range combinedScores {
		scoredDocs = append(scoredDocs, ScoredDoc{DocID: docID, Scores: scores})
	}

	// sort by HybridScore (desc)
	sort.Slice(scoredDocs, func(i, j int) bool {
		return scoredDocs[i].Scores.HybridScore > scoredDocs[j].Scores.HybridScore
	})

	if limit < len(scoredDocs) {
		scoredDocs = scoredDocs[:limit]
	}

	results := make([]WeightedSearchResult, len(scoredDocs))

	for i, d := range scoredDocs {
		doc := hs.Css.Documents[d.DocID]
		results[i] = WeightedSearchResult{
			DocID:         d.DocID,
			Title:         doc.Title,
			Description:   doc.Description,
			HybridScore:   d.Scores.HybridScore,
			KeywordScore:  d.Scores.KeywordScore,
			SemanticScore: d.Scores.SemanticScore,
		}
	}

	return results, nil

}

func (hs *HybridSearch) RRFSearch(query string, k int, limit int) ([]RRFSearchResult, error) {
	searchLimit := min(limit*500, len(hs.Css.Documents))

	// keyword search
	keywordResults, err := hs.bm25Search(query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform bm25Search: %v\n", err)
	}
	// semantic search
	semanticResults, err := hs.Css.SearchChunked(query, searchLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform SearchChunked: %v\n", err)
	}

	// combine results from both keyword and semantic
	// map of docId -> rank info + score
	combined := make(map[int]*CombinedRanks)

	// Fill keyword ranks (+score)
	for rank, result := range keywordResults {
		combined[result.DocID] = &CombinedRanks{
			KeywordRank:  rank,
			SemanticRank: -1,
			RRFScore:     CalcRRFScore(rank, k),
		}
	}

	// Fill semantic ranks (+score)
	for rank, result := range semanticResults {
		cr, ok := combined[result.DocID]
		if !ok {
			combined[result.DocID] = &CombinedRanks{
				KeywordRank:  -1,
				SemanticRank: rank,
				RRFScore:     CalcRRFScore(rank, k),
			}
		} else {
			cr.SemanticRank = rank
			cr.RRFScore = cr.RRFScore + CalcRRFScore(rank, k)
		}
	}

	// map -> slice for sorting
	rankedDocs := make([]RankedDoc, 0, len(combined))
	for docID, ranksInfo := range combined {
		rankedDocs = append(rankedDocs, RankedDoc{DocID: docID, Ranks: ranksInfo})
	}

	// sort by RRFScore (desc)
	sort.Slice(rankedDocs, func(i, j int) bool {
		return rankedDocs[i].Ranks.RRFScore > rankedDocs[j].Ranks.RRFScore
	})

	if limit < len(rankedDocs) {
		rankedDocs = rankedDocs[:limit]
	}

	results := make([]RRFSearchResult, len(rankedDocs))

	for i, d := range rankedDocs {
		doc := hs.Css.Documents[d.DocID]
		results[i] = RRFSearchResult{
			DocID:        d.DocID,
			Title:        doc.Title,
			Description:  doc.Description,
			RRFScore:     d.Ranks.RRFScore,
			KeywordRank:  d.Ranks.KeywordRank,
			SemanticRank: d.Ranks.SemanticRank,
		}
	}

	return results, nil
}

func Normalize(inputs []float64) []float64 {
	if len(inputs) == 0 {
		return []float64{}
	}

	minVal := slices.Min(inputs)
	maxVal := slices.Max(inputs)
	maxMinDiff := maxVal - minVal

	results := make([]float64, len(inputs))
	if maxMinDiff == 0 {
		for i := range results {
			results[i] = 1.0
		}
	} else {
		for i, v := range inputs {
			score := (v - minVal) / maxMinDiff
			results[i] = score
		}
	}

	return results
}

// alpha (or "α") is just a constant that we can use to dynamically control
// the weighting between the two scores

// Query Type	  	Example	          Chosen Alpha	  Reason
// Exact match	  "The Revenant"	  0.8	            Title search needs keywords
// Conceptual	  	"family movies"	  0.2	            Meaning matters more
// Mixed	        "2015 comedies"	  0.5	            Both year AND concept

func HybridScore(bm25Score float64, semanticScore float64, alpha float64) float64 {
	return alpha*bm25Score + (1-alpha)*semanticScore
}

// This is why it's so important to tune your search system's constants based on
// the types of data and queries you're working with in your application! It's not
// a one-size-fits-all solution, but building configurability into your system
// allows you to adjust it as needed.

func CalcRRFScore(rank int, k int) float64 {
	return 1 / float64(k+rank)
}
