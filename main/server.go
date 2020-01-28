package main

import (
	. "github.com/gorilla/mux"
	"github.com/rs/cors"
	. "kanban-distributed-system/message-queue"
	. "kanban-distributed-system/migration"
	. "kanban-distributed-system/request-handlers"
	. "kanban-distributed-system/utility"
	"log"
	"math/rand"
	. "net/http"
	"strconv"
	"time"
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

	rand.Seed(time.Now().UnixNano())

	port := 10000 + rand.Intn(1000)

	PORT = ":" + strconv.FormatInt(int64(port), 10)

	connection := CreateConnection()

	writeChan := CreateChannel(connection, PORT)

	channel := CreateChannel(connection, PORT)

	go ConsumeMessage(channel, PORT)

	RegisterServer(PORT)

	go UpdateMe(writeChan, PORT)

	log.Fatal(ListenAndServe(PORT, cors.Default().Handler(routes())))
}
