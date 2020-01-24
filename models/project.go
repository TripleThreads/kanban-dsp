package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	. "kanban-distributed-system/config"
	. "net/http"
)

type Project struct {
	gorm.Model
	Title string `json:"title"`
	Tasks []Task
}

// get list of project
func GetProjects(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	var projects []Project
	db.Find(&projects)
	json.NewEncoder(w).Encode(projects)
	_ = db.Close()
}

// get single project
func GetProject(w ResponseWriter, r *Request) {
	var project Project
	var tasks []Task
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	db.Where("id = ?", id).Find(&project)
	db.Where("project_id = ?", project.ID).Find(&tasks)
	project.Tasks = tasks
	fmt.Println(tasks)
	json.NewEncoder(w).Encode(project)
	_ = db.Close()
}

// new project
func CreateProject(w ResponseWriter, r *Request) {
	var project Project
	db := ConnectToMysql()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || json.Unmarshal(body, &project) != nil {
		Error(w, err.Error(), 500)
		return
	}

	db.Create(&project)
	json.NewEncoder(w).Encode(project)
	_ = db.Close()
}

// delete project
func DeleteProject(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	var project Project
	db.Where("id = ?", id).Find(&project)
	db.Delete(project)
	_ = db.Close()
}

// edit project
func UpdateProject(w ResponseWriter, r *Request) {
	db := ConnectToMysql()
	vars := mux.Vars(r)
	id := vars["id"]
	var project Project
	var newProj Project

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || json.Unmarshal(body, &newProj) != nil {
		Error(w, err.Error(), 500)
		return
	}

	db.Where("id = ?", id).Find(&project)

	project.Title = newProj.Title
	db.Save(project)
	_ = db.Close()
}
