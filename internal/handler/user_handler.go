package handler

import (
	"justn0w-bot/internal/request"
	"justn0w-bot/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (u UserHandler) Register(c *gin.Context) {
	registerRequest := &request.UserRequest{}
	err := c.ShouldBindJSON(registerRequest)
	if err != nil {
		log.Fatalln("参数解析错误")
	}

	userService := service.UserService{}
	err = userService.Register(registerRequest.Name, registerRequest.Password)
	if err != nil {
		ReturnError(c, 500, err.Error(), nil)
		return
	}
	ReturnSuccess(c, 200, "注册成功", nil)
}

func (u UserHandler) Login(context *gin.Context) {

}
