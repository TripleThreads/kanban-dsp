package main

import (
	. "github.com/gorilla/mux"
	. "kanban-distributed-system/migration"
	. "kanban-distributed-system/request-handlers"
	. "kanban-distributed-system/utility"
	"log"
	. "net/http"
)

func routes() *Router {
	router := NewRouter()

	// models routes
	router.HandleFunc("/projects/all", GetProjectsHandler).Methods("GET")
	router.HandleFunc("/projects/{ID}", GetProjectHandler).Methods("GET")
	router.HandleFunc("/projects/create", CreateProjectHandler).Methods("POST")
	router.HandleFunc("/projects/{ID}", DeleteProjectHandler).Methods("DELETE")
	router.HandleFunc("/projects/{ID}", UpdateProjectHandler).Methods("PUT")

	// tasks routes
	router.HandleFunc("/tasks/all", GetTasksHandler).Methods("GET")
	router.HandleFunc("/tasks/{ID}", GetTaskHandler).Methods("GET")
	router.HandleFunc("/tasks/create", CreateTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/{ID}", DeleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tasks/{ID}", UpdateTaskHandler).Methods("PUT")
	return router
}

func main() {
	InitialMigration()
	str := RegisterServer()

	log.Fatal(ListenAndServe(str, routes()))
}
