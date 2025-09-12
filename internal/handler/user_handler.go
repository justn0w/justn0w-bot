package handler

import (
	"justn0w-bot/internal/request"
	"justn0w-bot/internal/response"
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

	if registerRequest.Name == "" || registerRequest.Password == "" {
		response.ReturnError(c, 400, "参数错误", "")
		return
	}

	userService := service.UserService{}
	err = userService.Register(registerRequest.Name, registerRequest.Password)
	if err != nil {
		response.ReturnError(c, 500, err.Error(), nil)
		return
	}
	response.ReturnSuccess(c, 200, "注册成功", nil)
}

func (u UserHandler) Login(c *gin.Context) {
	registerRequest := &request.UserRequest{}
	err := c.ShouldBindJSON(registerRequest)
	if err != nil {
		log.Fatalln("参数解析错误")
	}
	if registerRequest.Name == "" || registerRequest.Password == "" {
		response.ReturnError(c, 400, "参数错误", "")
		return
	}

	userService := service.UserService{}
	loginResponse, err := userService.Login(registerRequest.Name, registerRequest.Password)
	if err != nil {
		response.ReturnError(c, 400, err.Error(), loginResponse)
		return
	}
	response.ReturnSuccess(c, 200, "登录成功", loginResponse)
}
