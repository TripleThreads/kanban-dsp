package migration

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/config"
	. "kanban-distributed-system/models"
)

func InitialMigration() {
	db := ConnectToMysql()
	db.AutoMigrate(Project{}, Task{})
	_ = db.Close()
}
