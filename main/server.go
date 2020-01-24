package main

import (
	. "github.com/gorilla/mux"
	. "kanban-distributed-system/migration"
	. "kanban-distributed-system/models"
	"log"
	. "net/http"
)

func routes() *Router {
	router := NewRouter()

	// models routes
	router.HandleFunc("/projects/all", GetProjects).Methods("GET")
	router.HandleFunc("/projects/{ID}", GetProject).Methods("GET")
	router.HandleFunc("/projects/create", CreateProject).Methods("POST")
	router.HandleFunc("/projects/{ID}/delete", DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{ID}/edit", UpdateProject).Methods("PUT")

	// tasks routes
	router.HandleFunc("/tasks/all", GetTask).Methods("GET")
	router.HandleFunc("/tasks/{ID}", GetTask).Methods("GET")
	router.HandleFunc("/tasks/create", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{ID}/delete", DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{ID}/edit", UpdateTask).Methods("PUT")
	return router
}

func main() {
	InitialMigration()

	log.Fatal(ListenAndServe(":8000", routes()))
}
