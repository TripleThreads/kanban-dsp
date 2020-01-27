package migration

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/commons"
	. "kanban-distributed-system/config"
	. "kanban-distributed-system/models"
)

func InitialMigration() {
	db := ConnectToMysql()
	db.AutoMigrate(Project{}, Task{}, Operation{})
	_ = db.Close()
}
