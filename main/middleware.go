package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gorilla/mux"
	"github.com/streadway/amqp"
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

var servers []string

func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

var lastServer = 0

func AddServer(connection *amqp.Connection, port string) {
	println("Registering ", port)

	servers = append(servers, port)

	msg, err := json.Marshal(servers)

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
	for i, port := range servers {
		channel := CreateChannel(connection, port)
		PublishMessage(channel, msg, servers[i])
		_ = channel.Close()
	}
}

func handleRequest(connection *amqp.Connection) *Router {
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
			AddServer(connection, msg.Port)
		}

	}).Methods("POST")

	router.HandleFunc("/port", func(writer ResponseWriter, request *Request) {
		json.NewEncoder(writer).Encode(lastServer)
		lastServer++
	}).Methods("GET")
	return router
}
func main() {
	connection := CreateConnection() // message queue
	println("hold on..")
	log.Fatal(ListenAndServe(":9865", handleRequest(connection)))
}
