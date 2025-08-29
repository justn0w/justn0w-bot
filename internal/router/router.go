package router

import (
	"justn0w-bot/internal/handler"

	"github.com/gin-gonic/gin"
)

var (
	doChatHandler handler.ChatHandler
)

func Router() *gin.Engine {
	r := gin.Default()

	{
		chatGroup := r.Group("/chat")
		chatGroup.POST("/generate", doChatHandler.Generate)
		chatGroup.POST("/generate/stream", doChatHandler.GenerateStream)
	}
	return r
}
