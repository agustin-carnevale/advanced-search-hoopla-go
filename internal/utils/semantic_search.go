package utils

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type SemanticSearch struct {
	Model  string
	client *api.Client
}

func NewSemanticSearch(modelName string) (*SemanticSearch, error) {
	client, err := api.ClientFromEnvironment()
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
	resp, err := ss.client.Embeddings(context.Background(), &api.EmbeddingRequest{
		Model:  ss.Model,
		Prompt: "test",
	})
	if err != nil {
		return err
	}

	fmt.Println("Vector dimensions:", len(resp.Embedding))
	return nil
}
