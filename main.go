package main

import (
	"log"
	"recrem/config/app"
)

func main() {
	engine := app.InitApp() // 初始化

	err := engine.Run(":" + setting.Config.Server.Port) // 运行
	if err != nil {
		log.Logger.Sugar().Panic("项目启动失败: ", err.Error())
	}
}
