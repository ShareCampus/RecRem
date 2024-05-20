package app

import (
	"log"
	"recrem/config/db"
	"recrem/config/migrate"
	"recrem/config/setting"
	logger "recrem/log"
	"recrem/middlewares"
	"recrem/routers"
	"reflect"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

// InitApp 初始化 gin
func InitApp() *gin.Engine {
	// 加载配置
	s := setting.Setting{}
	s.InitSetting()
	s.InitLute()
	s.InitCache()
	db.InitDb()
	migrate.Migrate()
	gin.SetMode(setting.Config.Server.Mode)

	// 加载中间件
	router := gin.New()
	err := logger.InitLogger(
		setting.Config.Logger.FileName,
		setting.Config.Logger.Level,
		setting.Config.Logger.Format,
		setting.Config.Logger.MaxSize,
		setting.Config.Logger.MaxBackups,
		setting.Config.Logger.MaxAge,
	)
	if err != nil {
		log.Panicln("初始化日志失败：", err.Error())
	}
	// 限流中间件
	middlewares.InitBucket(time.Second*time.Duration(setting.Config.Server.LimitTime), setting.Config.Server.LimitCap)
	store := cookie.NewStore([]byte("secret-recrem-store"))
	// 开启中间件
	router.Use(middlewares.Logger(logger.Logger), middlewares.Recover(logger.Logger, true),
		middlewares.Limiter(), sessions.Sessions("mySession", store))
	// 配置表单校验
	uni := ut.New(zh.New())
	setting.Trans, _ = uni.GetTranslator("zh")
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = translations.RegisterDefaultTranslations(v, setting.Trans)
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := field.Tag.Get("label")
			return name
		})
	}
	// 加载路由
	apiRouter := routers.ApiRouter{}
	tmplRouter := routers.TmplRouter{}
	tmplRouter.InitTemplateRouter("", router)
	apiRouter.InitApiRouter("/api/v1", router)
	return router
}
