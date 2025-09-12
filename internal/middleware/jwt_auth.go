package middleware

import (
	"justn0w-bot/internal/response"
	"justn0w-bot/internal/service"
	"justn0w-bot/pkg/rescode"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ReturnFailedWithErrorCode(c, rescode.TokenInvalid)
			c.Abort()
			return
		}

		splits := strings.SplitN(authHeader, " ", 2)
		if len(splits) != 2 || splits[0] != "Bearer" {
			response.ReturnFailedWithErrorCode(c, rescode.TokenAuthFormat)
			c.Abort()
			return
		}

		// 验证token
		claims, err := service.ParseToken(splits[1])
		if err != nil {
			response.ReturnFailedWithErrorCode(c, rescode.TokenInvalid)
			c.Abort()
			return
		}

		c.Set("userName", claims.UserName)
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
