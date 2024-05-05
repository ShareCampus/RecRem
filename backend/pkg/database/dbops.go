package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"os"
)

func SetupDB(database *gorm.DB) *gorm.DB {
	if database != nil {
		return database
	}
	var (
		user         = os.Getenv("DB_USER")
		password     = os.Getenv("DB_PASSWORD")
		address      = os.Getenv("DB_ADDRESS")
		port         = os.Getenv("DB_PORT")
		databaseName = os.Getenv("DB_NAME")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", user, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(glog.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
