package router

import (
	"justn0w-bot/internal/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	{
		chatGroup := r.Group("/chat")
		chatGroup.POST("/generate", handler.ChatHandler{}.Generate)
		chatGroup.POST("/generate/stream", handler.ChatHandler{}.GenerateStream)
	}

	{
		chatGroup := r.Group("/rag")
		chatGroup.POST("/file/upload", handler.RagHandler{}.UploadFile)
	}
	return r
}
