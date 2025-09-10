package config

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitMySQL() {
	dbClient, err := gorm.Open(mysql.Open(viper.GetString("mysql.url")), &gorm.Config{})
	if err != nil {
		panic("init mysql error")
	}
	log.Println("init mysql success")
	Db = dbClient
}
