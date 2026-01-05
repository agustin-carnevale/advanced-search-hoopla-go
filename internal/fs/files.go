package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agustin-carnevale/advanced-search-hoopla-go/internal/model"
)

const (
	DefaultSearchLimit = 5
)

var (
	ProjectRoot, _ = getProjectRoot()

	DataPath          = filepath.Join(ProjectRoot, "data", "movies.json")
	StopWordsPath     = filepath.Join(ProjectRoot, "data", "stopwords.txt")
	GoldenDatasetPath = filepath.Join(ProjectRoot, "data", "golden_dataset.json")
	CacheDir          = filepath.Join(ProjectRoot, "cache")
	IndexPath         = filepath.Join(CacheDir, "index.gob")
	EmbeddingsPath    = filepath.Join(CacheDir, "movie_embeddings.gob")

	ChunksEmbeddingsPath = filepath.Join(CacheDir, "chunks_embeddings.gob")
	ChunksMetadataPath   = filepath.Join(CacheDir, "chunks_metadata.json")
)

// getProjectRoot walks up until it finds go.mod (project base)
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // reached filesystem root
			return "", fmt.Errorf("project root not found (go.mod missing)")
		}
		dir = parent
	}
}

func LoadMovies() ([]model.Movie, error) {
	file, err := os.ReadFile(DataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", DataPath, err)
	}

	var data struct {
		Movies []model.Movie `json:"movies"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed parsing movies.json: %w", err)
	}

	return data.Movies, nil
}

func LoadStopWords() (map[string]struct{}, error) {
	raw, err := os.ReadFile(StopWordsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", StopWordsPath, err)
	}

	stopwords := make(map[string]struct{})

	// Split into lines and normalize
	lines := strings.Split(string(raw), "\n")

	for _, line := range lines {
		word := strings.TrimSpace(line)
		if word != "" {
			stopwords[word] = struct{}{}
		}
	}

	return stopwords, nil
}

// func splitLines(input string) []string {
// 	return filepath.SplitList(input) // handles \n and platform separators
// }

// func trim(s string) string {
// 	return os.ExpandEnv(s) // lazy but removes newlines & spaces
// }

func LoadGoldenDataset() ([]model.TestCase, error) {
	file, err := os.ReadFile(GoldenDatasetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", GoldenDatasetPath, err)
	}

	var data struct {
		TestCases []model.TestCase `json:"test_cases"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed parsing %s: %w", GoldenDatasetPath, err)
	}

	return data.TestCases, nil
}
