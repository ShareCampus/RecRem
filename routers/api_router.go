package routers

import (
	"recrem/handlers/api"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (a *ApiRouter) InitApiRouter(rootPath string, router *gin.Engine) {
	authHandler := api.AuthHandler{}
	userHandler := api.UserHandler{}
	fileHandler := api.FileHandler{}
	queryHandler := api.QueryHandler{}

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

	userApiRouter := router.Group(rootPath)
	{
		userApiRouter.GET("/all_users", userHandler.GetAllUsers)
		// userApiRouter.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		// userApiRouter.PUT("/users/pwd", middlewares.JWTAuth(), userHandler.UpdateUserPwd)
		// userApiRouter.PUT("/users", middlewares.JWTAuth(), userHandler.UpdateUser)
		// userApiRouter.PUT("/users/pwd", middlewares.JWTAuth(), userHandler.UpdateUserPwd)
	}

	fileApiRouter := router.Group(rootPath)
	{
		fileApiRouter.POST("/file/upload", fileHandler.UploadFile)
		fileApiRouter.DELETE("/file/delete", fileHandler.DeleteFile)
	}

	queryApiRouter := router.Group(rootPath)
	{
		queryApiRouter.POST("query/question", queryHandler.QueryByQuestion)
		queryApiRouter.POST("query/file", queryHandler.QueryByFile)
	}
}
