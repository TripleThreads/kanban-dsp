package projects

import (
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Task struct {
	gorm.Model
	Title  string        `json:"title"`
	Period time.Duration `json:"time"`
}

// get list of tasks
func getTasks(w http.ResponseWriter, r *http.Request) {

}

// get single task
func getTask(w http.ResponseWriter, r *http.Request) {

}

// new task
func createTask(w http.ResponseWriter, r *http.Request) {

}

// delete task
func deleteTask(w http.ResponseWriter, r *http.Request) {

}

// edit task
func editTask(w http.ResponseWriter, r *http.Request) {

}
