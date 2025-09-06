package service

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/spf13/viper"
)

func buildEmbedding(ctx context.Context) embedding.Embedder {
	baseURL := viper.GetString("ollama.base_url")
	model := viper.GetString("ollama.embedding_model")

	embedder, _ := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: baseURL,
		Model:   model,
		Timeout: 10 * time.Second,
	})

	return embedder
}
