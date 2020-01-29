package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/config"
	"strconv"
	"time"
)

type Task struct {
	gorm.Model
	Title     string `json:"title"`
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
func CreateTask(task Task, timestamp time.Time) []byte {
	db := ConnectToMysql()
	db.Create(&task)
	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "CREATE", "TASK", timestamp)
	return msg
}

// delete task
func DeleteTask(id string, timestamp time.Time) []byte {
	var task Task
	db := ConnectToMysql()
	db.Where("id = ?", id).Find(&task)
	db.Delete(task)

	ID, _ := strconv.Atoi(id)
	task.ID = uint(ID)

	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "DELETE", "TASK", timestamp)
	return msg
}

// edit task
func UpdateTask(id string, task Task, timestamp time.Time) []byte {
	var tk Task

	db := ConnectToMysql()
	db.Where("id = ?", id).Find(&tk)
	db.Model(&tk).Updates(&task)

	ID, _ := strconv.Atoi(id)
	task.ID = uint(ID)

	_ = db.Close()
	msg, err := json.Marshal(task)
	checkError(err)
	LogOperation(msg, "UPDATE", "TASK", timestamp)
	return msg
}
