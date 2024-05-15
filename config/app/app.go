package app

import (
	"recrem/config/setting"
	"recrem/routers"

	"github.com/gin-gonic/gin"
)

// InitApp 初始化 gin
func InitApp() *gin.Engine {
	// 加载配置
	s := setting.Setting{}
	s.InitSetting()
	s.InitLute()
	s.InitCache()
	// db.InitDb()
	// migrate.Migrate()
	gin.SetMode(setting.Config.Server.Mode)

	// 加载中间件
	router := gin.New()
	// 加载路由
	apiRouter := routers.ApiRouter{}
	// tmplRouter := routers.TmplRouter{}
	// tmplRouter.InitTemplateRouter("", router)
	apiRouter.InitApiRouter("/api/v1", router)

	return router
}
