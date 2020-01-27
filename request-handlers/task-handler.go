package request_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	. "kanban-distributed-system/message-queue"
	. "kanban-distributed-system/models"
	. "kanban-distributed-system/utility"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		println(err)
	}
}

// Get list of tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := GetTasks()
	json.NewEncoder(w).Encode(tasks)
}

// Get single task
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	task := GetTask(id)
	json.NewEncoder(w).Encode(task)
}

// new task
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || json.Unmarshal(body, &task) != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	msg := CreateTask(task)
	json.NewEncoder(w).Encode(task)
	PropagateUpdate(msg, PORT)
}

// delete task
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	msg := DeleteTask(id)
	json.NewEncoder(w).Encode("0k")
	PropagateUpdate(msg, PORT)
}

// edit task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var task Task
	id := vars["ID"]
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	fmt.Println(task)
	msg := UpdateTask(id, task)
	PropagateUpdate(msg, PORT)
}
