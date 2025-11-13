package utils

import (
	"context"
	"fmt"

	ollama "github.com/ollama/ollama/api"
)

type SemanticSearch struct {
	Model  string
	client *ollama.Client
}

func NewSemanticSearch(modelName string) (*SemanticSearch, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &SemanticSearch{
		Model:  modelName,
		client: client,
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
