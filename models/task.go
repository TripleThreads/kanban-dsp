package models

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/config"
	"time"
)

type Task struct {
	gorm.Model
	Title     string    `json:"title"`
	Period    time.Time `json:"time"`
	ProjectId uint
	Project   Project `gorm:"association_foreignkey:Id"`
	Stage     uint8
}

// Get list of tasks
func GetTasks() []Task {
	db := ConnectToMysql()
	var tasks []Task
	db.Find(&tasks)
	_ = db.Close()
	return tasks
}

// Get single task
func GetTask(id string) Task {
	db := ConnectToMysql()
	var task Task
	db.Where("id = ?", id).Find(&task)
	_ = db.Close()
	return task
}

// new task
func CreateTask(task Task) []byte {
	db := ConnectToMysql()
	db.Create(&task)
	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "CREATE", "TASK")
	return msg
}

// delete task
func DeleteTask(id string) []byte {
	db := ConnectToMysql()
	var task Task
	db.Where("id = ?", id).Find(&task)
	db.Delete(task)
	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "DELETE", "TASK")
	return msg
}

// edit task
func UpdateTask(id string, task Task) []byte {
	db := ConnectToMysql()
	var tk Task
	db.Where("id = ?", id).Find(&tk)
	fmt.Println(id)
	fmt.Println(task)
	db.Model(&tk).Updates(&task)
	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "CREATE", "TASK")
	return msg
}
