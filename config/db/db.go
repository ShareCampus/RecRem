package db

import (
	"fmt"
	"log"
	"net/url"
	"recrem/config/setting"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // mysql驱动
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // gorm mysql
)

var Db *gorm.DB

func getDataSource() string {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		setting.Config.Database.UserName,
		setting.Config.Database.Password,
		setting.Config.Database.Host,
		setting.Config.Database.Port,
		setting.Config.Database.Database,
		url.QueryEscape(setting.Config.Database.TimeZone), // 对时区进行 Url 编码
	)
	return dataSource
}

func InitDb() {
	var err error
	Db, err = gorm.Open("mysql", getDataSource())
	if err != nil {
		log.Panic("数据库连接错误：", err.Error())
	}
	Db.DB().SetMaxIdleConns(setting.Config.Database.MaxIdleConn)
	Db.DB().SetMaxOpenConns(setting.Config.Database.MaxOpenConn)
	if setting.Config.Server.Mode == gin.DebugMode {
		Db.LogMode(true)
	}
}
