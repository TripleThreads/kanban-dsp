package request_handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	. "kanban-distributed-system/message-queue"
	. "kanban-distributed-system/models"
	. "kanban-distributed-system/utility"
	. "net/http"
)

// new project
func CreateProjectHandler(w ResponseWriter, r *Request) {
	var project Project

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || json.Unmarshal(body, &project) != nil {
		Error(w, err.Error(), 500)
		return
	}
	msg := CreateProject(project)
	json.NewEncoder(w).Encode(project)
	PropagateUpdate(msg, PORT)
}

// get list of project
func GetProjectsHandler(w ResponseWriter, r *Request) {
	projects := GetProjects()
	json.NewEncoder(w).Encode(projects)
}

// get single project
func GetProjectHandler(w ResponseWriter, r *Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	project := GetProject(id)
	json.NewEncoder(w).Encode(project)
}

// delete project
func DeleteProjectHandler(w ResponseWriter, r *Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	msg := DeleteProject(id)
	json.NewEncoder(w).Encode("0k")
	PropagateUpdate(msg, PORT)
}

// edit project
func UpdateProjectHandler(w ResponseWriter, r *Request) {

	vars := mux.Vars(r)
	id := vars["ID"]
	var project Project

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || json.Unmarshal(body, &project) != nil {
		Error(w, err.Error(), 500)
		return
	}
	msg := UpdateProject(id, project)
	PropagateUpdate(msg, PORT)
}