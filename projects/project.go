package projects

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/config"
	. "net/http"
)

func InitialMigration() {
	db := ConnectToMysql()
	db.AutoMigrate(Project{})
	_ = db.Close()
}

type Project struct {
	gorm.Model
	Title string `json:"title"`
	Tasks *Task  `json:"tasks"`
}

// get list of tasks
func GetProjects(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	var projects []Project
	db.Find(&projects)
	json.NewEncoder(w).Encode(projects)
	_ = db.Close()
}

// get single task
func getProject(w ResponseWriter, r *Request) {

}

// new task
func CreateProject(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)

	title := vars["title"]

	db.Create(&Project{Title: title})
	db.Close()
}

// delete task
func DeleteProject(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	var project Project
	db.Where("id = ?", id).Find(&Project{})
	db.Delete(project)
	_ = db.Close()
}

// edit task
func UpdateProject(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	title := vars["title"]
	var project Project
	db.Where("id = ?", id).Find(&project)

	project.Title = title
	db.Save(project)
	_ = db.Close()
}
