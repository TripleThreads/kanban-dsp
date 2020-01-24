package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	. "kanban-distributed-system/config"
	"net/http"
	"time"
)

type Task struct {
	gorm.Model
	Title     string    `json:"title"`
	Period    time.Time `json:"time"`
	ProjectId uint
	Project   Project `gorm:"foreignkey:ProjectId"`
}

// Get list of tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	db := ConnectToMysql()
	var tasks []Task
	db.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
	_ = db.Close()
}

// Get single task
func GetTask(w http.ResponseWriter, r *http.Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	var task Task
	db.Where("id = ?", id).Find(&task)
	json.NewEncoder(w).Encode(task)
	_ = db.Close()
}

// new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	db := ConnectToMysql()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || json.Unmarshal(body, &task) != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	db.Create(&task)
	json.NewEncoder(w).Encode(task)
	_ = db.Close()
}

// delete task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	var task Task
	db.Where("id = ?", id).Find(task)
	db.Delete(task)
	_ = db.Close()
}

// edit task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	title := vars["title"]
	period := vars["period"]
	var task Task
	db.Where("id = ?", id).Find(&task)
	task.Title = title
	layout := "2006-01-02T15:04:05.000Z"

	// watch out for the following line might cause problem
	task.Period, _ = time.Parse(layout, period)
	_ = db.Close()
}
