package main

import (
	"fmt"
	"recrem/config/app"
	"recrem/config/setting"
)

// @title Gin Swagger
// @version 1.0
// @description recrem 开源搜索引擎 API 接口文档

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8088
func main() {
	engine := app.InitApp() // 初始化

	err := engine.Run(":" + setting.Config.Server.Port) // 运行
	fmt.Println("项目启动失败", err)
}
