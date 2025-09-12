package response

import (
	"justn0w-bot/pkg/rescode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}) {
	json := &CommonResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}, data interface{}) {
	json := &CommonResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, json)
}

func ReturnFailedWithErrorCode(c *gin.Context, code rescode.ErrorCode) {
	ReturnError(c, code.Code, code.Message, nil)
}
