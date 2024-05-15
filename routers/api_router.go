package routers

import (
	"recrem/handlers/api"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (a *ApiRouter) InitApiRouter(rootPath string, router *gin.Engine) {
	// 注册功能
	authHandler := api.AuthHandler{}
	authApiRouter := router.Group(rootPath)
	{
		authApiRouter.GET("/auth/health", authHandler.Health)
		authApiRouter.POST("/auth/register", authHandler.Register)
		authApiRouter.POST("/auth/login", authHandler.Login)

	}
}
