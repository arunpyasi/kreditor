package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/mdeheij/kreditor/config"
)

var (
	Database *gorm.DB
)

func init() {

	if config.C.DatabaseURI == "" {
		fmt.Println("Please specify MySQL database URI using database environment variable.")
		config.C.DatabaseURI = "root@tcp(127.0.0.1:3306)/kreditor"

	}

	var err error

	Database, err = gorm.Open("mysql", config.C.DatabaseURI+"?charset=utf8&parseTime=True&loc=Local")

	Database.Debug()
	Database.LogMode(true)

	if err != nil {
		panic("Failed to connect database:" + err.Error())
	}
	Database.AutoMigrate(&Debt{})
	Database.AutoMigrate(&Invoice{})
	Database.AutoMigrate(&User{})
	Database.AutoMigrate(&Contact{})

}
