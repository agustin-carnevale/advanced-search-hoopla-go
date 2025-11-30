package utils

import (
	"context"
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/fs"
	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"

	ollama "github.com/ollama/ollama/api"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type SimilarityScore struct {
	Score float64
	Movie model.Movie
}

type SemanticSearchResult struct {
	DocID       int
	Score       float64
	Title       string
	Description string
}

type SemanticSearch struct {
	Model       string
	client      *ollama.Client
	Documents   []model.Movie
	DocumentMap map[int]model.Movie
	Embeddings  [][]float64
}

func NewSemanticSearch(modelName string) (*SemanticSearch, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &SemanticSearch{
		Model:       modelName,
		client:      client,
		DocumentMap: make(map[int]model.Movie),
	}, nil
}

func (ss *SemanticSearch) VerifyModel() error {
	fmt.Println("Model loaded:", ss.Model)
	resp, err := ss.client.Embeddings(context.Background(), &ollama.EmbeddingRequest{
		Model:  ss.Model,
		Prompt: "test",
	})
	if err != nil {
		return err
	}

	fmt.Println("Vector dimensions:", len(resp.Embedding))
	return nil
}

func (ss *SemanticSearch) EmbedText(text string) ([]float64, error) {
	resp, err := ss.client.Embeddings(context.Background(), &ollama.EmbeddingRequest{
		Model:  ss.Model,
		Prompt: text,
	})
	if err != nil {
		return nil, err
	}

	return resp.Embedding, nil
}

func (ss *SemanticSearch) BuildEmbeddings() ([][]float64, error) {
	fmt.Println("🔄 Building embeddings…")

	// Build strings: "title: description"
	strings := make([]string, len(ss.Documents))
	for i, doc := range ss.Documents {
		ss.DocumentMap[doc.ID] = doc
		strings[i] = fmt.Sprintf("%s: %s", doc.Title, doc.Description)
	}

	// Generate embeddings one by one (Ollama does not batch today)
	// embeddings := make([][]float64, 0, len(strings))
	// for i, text := range strings {
	// 	fmt.Println("Creating embedding for doc:", i)
	// 	resp, err := ss.client.Embeddings(context.Background(), &api.EmbeddingRequest{
	// 		Model:  ss.Model,
	// 		Prompt: text,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	embeddings = append(embeddings, resp.Embedding)
	// }
	// ss.Embeddings = embeddings

	// Parallel embeddings creation
	embeddings, err := ss.createEmbeddingsParallel(strings)
	if err != nil {
		return nil, err
	}
	ss.Embeddings = embeddings

	// Save to disk
	if err := ss.saveEmbeddings(); err != nil {
		return nil, err
	}

	fmt.Println("✅ Embeddings built and saved.")
	return embeddings, nil
}

func (ss *SemanticSearch) saveEmbeddings() error {
	file, err := os.Create(fs.EmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(ss.Embeddings)
}

func (ss *SemanticSearch) loadEmbeddings() error {
	file, err := os.Open(fs.EmbeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&ss.Embeddings)
}

func (ss *SemanticSearch) createEmbeddingsParallel(strings []string) ([][]float64, error) {
	docCount := len(strings)
	workerCount := runtime.NumCPU()

	fmt.Println("Starting embedding generation")
	fmt.Printf("Documents: %d | Workers: %d\n", docCount, workerCount)

	// Output slice (pre-allocate)
	embeddings := make([][]float64, docCount)

	// Channels
	jobs := make(chan int, docCount)
	errChan := make(chan error, 1) // allow only 1 error to signal shutdown
	wg := sync.WaitGroup{}

	// === Progress bar setup ===
	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(int64(docCount),
		mpb.PrependDecorators(
			decor.Name("Embedding: "),
			decor.CountersNoUnit("%d/%d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
		),
	)

	// Worker function
	worker := func(id int) {
		defer wg.Done()

		for idx := range jobs {
			text := strings[idx]

			// fmt.Printf("Worker %d processing doc number: %d\n", id, idx+1)

			resp, err := ss.client.Embeddings(context.Background(), &ollama.EmbeddingRequest{
				Model:  ss.Model,
				Prompt: text,
			})
			if err != nil {
				// first error wins, avoids blocking
				select {
				case errChan <- fmt.Errorf("embedding error on doc %d: %w", idx, err):
				default:
				}
				return
			}

			embeddings[idx] = resp.Embedding

			// Increment progress bar (it handles it thread-safely internally)
			bar.Increment()
		}
	}

	// Start workers
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(i + 1)
	}

	// Push jobs
	for i := 0; i < docCount; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for workers
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// If any worker errors, stop early
	select {
	case err := <-errChan:
		return nil, err
	case <-done:
	}

	// Must wait for mpb to flush + close
	p.Wait()

	return embeddings, nil
}

func (ss *SemanticSearch) LoadOrCreateEmbeddings(docs []model.Movie) ([][]float64, error) {
	ss.Documents = docs
	for _, doc := range docs {
		ss.DocumentMap[doc.ID] = doc
	}

	// If file exists → try loading
	if _, err := os.Stat(fs.EmbeddingsPath); err == nil {
		if err := ss.loadEmbeddings(); err == nil {
			// Verify vector count matches document count
			if len(ss.Embeddings) == len(ss.Documents) {
				fmt.Println("📂 Loaded embeddings from disk.")
				return ss.Embeddings, nil
			}
		}
	}

	// Else create them
	return ss.BuildEmbeddings()
}

func (ss *SemanticSearch) Search(query string, limit int) ([]SemanticSearchResult, error) {
	if len(ss.Embeddings) == 0 || len(ss.Embeddings) != len(ss.Documents) {
		return nil, fmt.Errorf("No embeddings loaded. Call `load_or_create_embeddings` first.")
	}

	query_embedding, err := ss.EmbedText(query)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create embedding of the query: %v\n", err)
	}

	similarities := make([]SimilarityScore, 0, len(ss.Embeddings))

	for i, doc_embedding := range ss.Embeddings {
		similarity_score := CosineSimilarity(query_embedding, doc_embedding)
		similarities = append(similarities, SimilarityScore{
			Score: similarity_score,
			Movie: ss.Documents[i],
		})
	}

	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})

	results := make([]SemanticSearchResult, 0, limit)
	for _, item := range similarities[:limit] {
		results = append(results, SemanticSearchResult{
			DocID:       item.Movie.ID,
			Score:       item.Score,
			Title:       item.Movie.Title,
			Description: item.Movie.Description,
		})
	}

	return results, nil
}

func CosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return 0.0 // or panic, but returning 0 is safer
	}

	var dot, norm1, norm2 float64
	for i := range vec1 {
		dot += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dot / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func SemanticChunk(text string, maxChunkSize int, overlap int) []string {
	// remove leading and trailing whitespace
	text = strings.TrimSpace(text)

	// TODO: this is removing the punctuation (add it back?)
	sentenceRegex := regexp.MustCompile(`([.!?])\s+`)
	rawSentences := sentenceRegex.Split(text, -1)

	sentences := make([]string, 0)
	for _, s := range rawSentences {
		// remove leading and trailing whitespace
		sentence := strings.TrimSpace(s)
		if len(sentence) > 0 {
			sentences = append(sentences, sentence)
		}
	}

	if len(sentences) == 0 {
		return []string{}
	}

	var chunks []string
	start := 0
	end := min(maxChunkSize, len(sentences))
	c := strings.Join(sentences[start:end], " ")
	chunks = append(chunks, c)

	start = min(maxChunkSize-overlap, len(sentences))
	for start+overlap < len(sentences) {
		end = min(start+maxChunkSize, len(sentences))
		c = strings.Join(sentences[start:end], " ")
		chunks = append(chunks, c)
		start = end - overlap
	}

	return chunks
}
