package handler

import (
	"justn0w-bot/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
}

type ChatRequest struct {
	Question string `json:"question"`
}

func (t ChatHandler) Generate(c *gin.Context) {
	question := c.DefaultPostForm("question", "")

	chatService := service.ChatService{}
	res, err := chatService.Generate(question)
	if err != nil {
		panic("失败")
	}
	ReturnSuccess(c, 200, "success", res)
}

func (t ChatHandler) GenerateStream(c *gin.Context) {
	chatRequest := ChatRequest{}
	err := c.ShouldBind(&chatRequest)
	if err != nil {
		ReturnError(c, 400, "参数错误", "")
		return
	}

	chatService := service.ChatService{}
	question := c.DefaultPostForm("question", "")
	chatService.GenerateStream(c, question)
}
