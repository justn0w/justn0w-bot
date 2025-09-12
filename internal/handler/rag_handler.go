package handler

import (
	"context"
	"justn0w-bot/internal/response"
	"justn0w-bot/internal/service"

	"github.com/gin-gonic/gin"
)

type RagHandler struct {
}

func (h RagHandler) UploadFile(c *gin.Context) {

	ctx := context.Background()

	// 读取文件内容
	file, err := c.FormFile("file")
	if err != nil {
		response.ReturnError(c, 500, "上传文件失败", err.Error())
		return
	}

	// 文本向量化
	err = service.RagService{}.Vectorize(file, ctx)
	if err != nil {
		response.ReturnError(c, 500, "向量化失败", err.Error())
	}

}
