package handler

import (
	"justn0w-bot/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	chatService service.ChatService
)

type ChatHandler struct{}

func (t *ChatHandler) DoChat(c *gin.Context) {
	question := c.DefaultPostForm("question", "")
	res, err := chatService.DoChat(question)
	if err != nil {
		panic("失败")
	}
	ReturnSuccess(c, 200, "success", res)
}
