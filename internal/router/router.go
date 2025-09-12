package router

import (
	"justn0w-bot/internal/handler"
	"justn0w-bot/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 在生产环境中应该指定具体的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	chatGroup := r.Group("/chat")
	{
		chatGroup.Use(middleware.JwtAuthMiddleWare())
		chatGroup.POST("/generate", handler.ChatHandler{}.Generate)
		chatGroup.POST("/generate/stream", handler.ChatHandler{}.GenerateStream)
	}

	ragGroup := r.Group("/rag")
	{
		ragGroup.POST("/file/upload", handler.RagHandler{}.UploadFile)
	}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", handler.UserHandler{}.Register)
		userGroup.POST("/login", handler.UserHandler{}.Login)
	}
	return r
}
