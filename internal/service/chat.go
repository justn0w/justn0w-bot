package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ChatService struct {
}

func (t ChatService) Generate(question string) (string, error) {
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

func (t ChatService) GenerateStream(c *gin.Context, question string) {
	//1 初始化context
	ctx := context.Background()

	//2 创建模型
	cm := buildModel(ctx)

	//3 对话内容
	messages := buildMessage(question)

	//4 拿到结果
	streamResult, _ := cm.Stream(ctx, messages)

	//5 流式输出
	handleStreamResponse(streamResult, c)
}

func handleStreamResponse(sr *schema.StreamReader[*schema.Message], c *gin.Context) {
	defer sr.Close()

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 用于客户端断开连接时退出循环
	clientGone := c.Writer.CloseNotify()

	for {
		select {
		case <-clientGone:
			// 客户端断开连接，退出
			log.Println("Client disconnected")
			c.SSEvent("close", "客户端断开连接")
			c.Writer.Flush()
			return
		default:
			content, err := sr.Recv()
			if err == io.EOF {
				c.SSEvent("close", "对话结束")
				c.Writer.Flush()
				return
			}
			//message事件和message的内容
			c.SSEvent("message", content)
			// 立即发送数据
			c.Writer.Flush()
		}
	}

}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
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
