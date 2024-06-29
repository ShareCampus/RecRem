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
	userHandler := api.UserHandler{}
	fileHandler := api.FileHandler{}
	authApiRouter := router.Group(rootPath)
	{
		authApiRouter.POST("/auth/login", authHandler.Login)
		authApiRouter.POST("/auth/register", authHandler.Register)
		// authApiRouter.GET("/auth/captcha", authHandler.CreateCaptcha)
		// authApiRouter.POST("/auth/pwd/forget", authHandler.ForgetPwd)
		// authApiRouter.POST("/auth/pwd/reset", authHandler.ResetPwd)
		// authApiRouter.GET("/auth/captcha", authHandler.CreateCaptcha)
		// authApiRouter.POST("/auth/pwd/forget", authHandler.ForgetPwd)
		// authApiRouter.POST("/auth/pwd/reset", authHandler.ResetPwd)
	}

	// 用户信息
	userApiRouter := router.Group(rootPath)
	{
		userApiRouter.GET("/all_users", userHandler.GetAllUsers)
		// userApiRouter.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		// userApiRouter.PUT("/users/pwd", middlewares.JWTAuth(), userHandler.UpdateUserPwd)
		// userApiRouter.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		// userApiRouter.PUT("/users/pwd", middlewares.JWTAuth(), userHandler.UpdateUserPwd)
	}
	// 上传文件，支持pdf, doc, md，txt文件解析和存储服务
	fileApiRouter := router.Group(rootPath)
	{
		fileApiRouter.POST("/file/upload", fileHandler.UploadFile)
		fileApiRouter.DELETE("/file/delete", fileHandler.DeleteFile)
	}

}
