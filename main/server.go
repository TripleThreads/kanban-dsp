package main

import (
	. "github.com/gorilla/mux"
	. "kanban-distributed-system/projects"
	"log"
	. "net/http"
)

func routes() *Router {
	router := NewRouter()
	router.HandleFunc("/projects", GetProjects).Methods("GET")
	router.HandleFunc("projects/{title}", CreateProject).Methods("POST")
	router.HandleFunc("/projects/{ID}", DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{ID}", UpdateProject).Methods("PUT")

	return router
}

func main() {
	InitialMigration()

	log.Fatal(ListenAndServe(":8000", routes()))
}
