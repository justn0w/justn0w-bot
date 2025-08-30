package service

import (
	"context"
	"justn0w-bot/config"
	"justn0w-bot/pkg/consts"
	"justn0w-bot/pkg/utils"
	"time"

	"log"
	"mime/multipart"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/cloudwego/eino-ext/components/indexer/redis"
)

type RagService struct {
}

// Vectorize 向量化
func (s RagService) Vectorize(uploadFile *multipart.FileHeader, ctx context.Context) error {
	log.Println(uploadFile.Filename)

	//// 1 读取文件的内容
	content, err := utils.GetFileContent(uploadFile)
	if err != nil {
		return err
	}
	//2 文本切割
	splitResults := splitText(content, ctx)

	//3 向量化和存储到redis库中
	saveVectors(splitResults, ctx)
	return nil

}

func saveVectors(docs []*schema.Document, ctx context.Context) {
	indexer, err := redis.NewIndexer(ctx, &redis.IndexerConfig{
		Client:    config.RedisClient,
		KeyPrefix: consts.RedisKeyPrefix,
		Embedding: buildEmbedding(ctx),
	})
	if err != nil {
		panic(err)
	}

	ids, _ := indexer.Store(ctx, docs)
	log.Printf("result: %v\n", ids)

}

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

func processEmbedding(docs []string, ctx context.Context) [][]float64 {
	baseURL := "http://127.0.0.1:11434"
	model := "nomic-embed-text"

	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		BaseURL: baseURL,
		Model:   model,
		Timeout: 10 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	vectors, err := embedder.EmbedStrings(ctx, docs)
	if err != nil {
		panic(err)
	}
	//log.Printf("vectors : %v", vectors)
	return vectors
}

func splitText(content string, ctx context.Context) []*schema.Document {

	//1 初始化分割器
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   1000,
		OverlapSize: 200,
		Separators:  []string{"\\n", "\\r\\n", "\\r", " ", "\\n\\n"},
		KeepType:    recursive.KeepTypeEnd,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return uuid.New().String()
		},
	})
	if err != nil {
		panic(err)
	}

	//2 准备要分割的文档
	docs := []*schema.Document{
		{
			Content: content,
		},
	}

	//3 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}
	log.Printf("results: %v\n", results)

	return results
}

//func splitText(content string, ctx context.Context) []string {
//
//	//1 初始化分割器
//	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
//		ChunkSize:   1000,
//		OverlapSize: 200,
//		Separators:  []string{"\\n", "\\r\\n", "\\r", " ", "\\n\\n"},
//		KeepType:    recursive.KeepTypeEnd,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	//2 准备要分割的文档
//	docs := []*schema.Document{
//		{
//			ID:      "doc1",
//			Content: content,
//		},
//	}
//
//	//3 执行分割
//	results, err := splitter.Transform(ctx, docs)
//
//	if err != nil {
//		panic(err)
//	}
//
//	resultContents := make([]string, len(results))
//	for i, result := range results {
//		resultContents[i] = result.Content
//	}
//	return resultContents
//}
