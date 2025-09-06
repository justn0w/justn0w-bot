package service

import (
	"context"
	"fmt"
	"justn0w-bot/config"
	"justn0w-bot/internal/request"
	"justn0w-bot/pkg/utils"
	"time"

	"log"
	"mime/multipart"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	milvus_retriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/spf13/viper"
)

var fields = []*entity.Field{
	{
		Name:     "id",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "256",
		},
		PrimaryKey: true,
	},
	{
		Name:     "vector", // 确保字段名匹配
		DataType: entity.FieldTypeBinaryVector,
		TypeParams: map[string]string{
			"dim": "24576",
		},
	},
	{
		Name:     "content",
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": "8192",
		},
	},
	{
		Name:     "metadata",
		DataType: entity.FieldTypeJSON,
	},
}

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
	emb := buildEmbedding(ctx)
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     config.MilvusClient,
		Embedding:  emb,
		Fields:     fields,
		Collection: viper.GetString("milvus.collection"),
	})
	if err != nil {
		log.Fatalf("Failed to create indexer: %v", err)
		return
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
		return
	}
	log.Printf("Store success, ids: %v", ids)

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
		ChunkSize:   600,
		OverlapSize: 50,
		Separators:  []string{"\n\n", "\n", "\r\n", "\r", ". ", "; ", " ", "\t"},
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

func handleRetrieval(ctx context.Context, chatRequest request.ChatRequest) []*schema.Document {
	emb := buildEmbedding(ctx)

	retriever, err := milvus_retriever.NewRetriever(ctx, &milvus_retriever.RetrieverConfig{
		Client:      config.MilvusClient,
		Collection:  viper.GetString("milvus.collection"),
		Partition:   nil,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
		},
		TopK:      1,
		Embedding: emb,
	})

	if err != nil {
		panic(err)
	}

	docs, err := retriever.Retrieve(ctx, chatRequest.Question)
	if err != nil {
		panic(err)
	}
	for i, doc := range docs {
		fmt.Printf("[%d] %s\n", i, doc.Content)
	}
	return docs
}
