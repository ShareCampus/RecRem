package routers

import (
	"recrem/handlers/api"

	"recrem/middlewares"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (a *ApiRouter) InitApiRouter(rootPath string, router *gin.Engine) {
	// 注册功能
	authHandler := api.AuthHandler{}
	userHandler := api.UserHandler{}
	fileHandler := api.FileHandler{}
	authApiRouter := router.Group(rootPath)
	{
		authApiRouter.POST("/auth/login", authHandler.Login)
		authApiRouter.POST("/auth/register", authHandler.Register)
		authApiRouter.GET("/auth/captcha", authHandler.CreateCaptcha)
		authApiRouter.POST("/auth/pwd/forget", authHandler.ForgetPwd)
		authApiRouter.POST("/auth/pwd/reset", authHandler.ResetPwd)
	}

	// 用户信息
	userApiRouter := router.Group(rootPath)
	{
		userApiRouter.GET("/all_users", userHandler.GetAllUsers)
		userApiRouter.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		userApiRouter.PUT("/users/pwd", middlewares.JWTAuth(), userHandler.UpdateUserPwd)
	}

	// 文件信息
	fileApiRouter := router.Group(rootPath)
	{
		fileApiRouter.POST("/upload_file", fileHandler.UploadFile) // 单文件上传
	}

}
