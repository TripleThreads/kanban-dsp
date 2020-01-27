package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gorilla/mux"
	"io/ioutil"
	. "kanban-distributed-system/commons"
	. "kanban-distributed-system/message-queue"
	"log"
	. "net/http"
	"time"
)

/*
* this program tracks nodes states
 */

var server []string

func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func AddServer(port string) {
	println("Registering ", port)
	connection := CreateConnection() // message queue

	server = append(server, port)
	msg, err := json.Marshal(server)

	var operations []Operation
	operations = append(operations, Operation{
		Sequence: time.Now(),
		Data:     msg,
	})

	message := Message{
		RequestType: "SERVERS",
		Port:        port,
		Operations:  operations,
	}
	msg, _ = json.Marshal(message)
	HandleError(err)
	for _, p := range server {
		channel := CreateChannel(connection, p)
		PublishMessage(channel, msg, p)
		_ = channel.Close()
	}

}

func handleRequest() *Router {
	router := NewRouter()
	router.HandleFunc("/", func(writer ResponseWriter, request *Request) {
		body, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		HandleError(err)
		var msg Message
		err = json.Unmarshal(body, &msg)
		HandleError(err)
		fmt.Println(msg)
		if msg.RequestType == "REG" { // let 1 be register new node
			AddServer(msg.Port)
		}

	}).Methods("POST")

	return router
}
func main() {
	println("hold on..")
	log.Fatal(ListenAndServe(":9865", handleRequest()))
}
