package main

import (
	"fmt"
	"recrem/config/app"
	"recrem/config/setting"
)

func main() {
	engine := app.InitApp() // 初始化

	err := engine.Run(":" + setting.Config.Server.Port) // 运行
	fmt.Println("项目启动失败", err)
}
