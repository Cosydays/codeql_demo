package client

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDb *gorm.DB
)

func InitMysql() {
	dsn := "root:root123@tcp(127.0.0.1:3306)/test_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	MysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
