package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func ConnectToMysql() *gorm.DB {
	DB, err := gorm.Open("mysql", "webuser:root@/kanban")
	if err != nil {
		fmt.Print(err.Error())
		panic("database connection failed")
	}
	return DB
}
