package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

type ChatService struct {
}

func (t *ChatService) DoChat(question string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loding .env file")
	}

	ctx := context.Background()

	cm := buildModel(ctx)
	messages := buildMessage(question)
	resp, err := cm.Generate(ctx, messages)
	if err != nil {
		log.Printf("Generate error: %v", err)
		return "", err
	}

	reasoning, ok := deepseek.GetReasoningContent(resp)
	if !ok {
		fmt.Printf("Unexpected: non-reasoning")
	} else {
		fmt.Printf("Resoning Content: %s\n", reasoning)
	}
	fmt.Printf("Assistant: %s\n", resp.Content)
	if resp.ResponseMeta != nil && resp.ResponseMeta.Usage != nil {
		fmt.Printf("Tokens used: %d (prompt) + %d (completion) = %d (total)\n",
			resp.ResponseMeta.Usage.PromptTokens,
			resp.ResponseMeta.Usage.CompletionTokens,
			resp.ResponseMeta.Usage.TotalTokens)
	}
	return resp.Content, err
}

func buildModel(ctx context.Context) *deepseek.ChatModel {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	// 创建 deepseek 模型
	cm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:    apiKey,
		Model:     "deepseek-reasoner",
		MaxTokens: 2000,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cm
}

func buildMessage(question string) []*schema.Message {
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "You are a helpful AI assistant. Be concise in your responses.",
		},
		{
			Role:    schema.User,
			Content: question,
		},
	}
	return messages
}
