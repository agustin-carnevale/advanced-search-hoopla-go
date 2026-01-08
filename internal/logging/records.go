package logging

type EnhancedQueryLog struct {
	EnhancementType string `json:"enhancement_type"`
	OriginalQuery   string `json:"original_query"`
	EnhancedQuery   string `json:"enhanced_query"`
}

type RRFCandidateLog struct {
	DocID        string  `json:"doc_id"`
	RRFScore     float64 `json:"rrf_score"`
	BM25Rank     int     `json:"bm25_rank"`
	SemanticRank int     `json:"semantic_rank"`
}

type FinalResultLog struct {
	DocID      string  `json:"doc_id"`
	FinalScore float64 `json:"final_score"`
	Position   int     `json:"position"`
}