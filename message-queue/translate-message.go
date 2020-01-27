package message_queue

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"kanban-distributed-system/commons"
	"kanban-distributed-system/models"
	"time"
)

// operation types
var DELETE = "DELETE"
var UPDATE = "UPDATE"
var CREATE = "CREATE"

// data types
var TASK = "TASK"
var PROJECT = "PROJECT"

// request types
var YOUR_UPDATE = "YOUR-UPDATE"
var UPDATE_ME = "UPDATE-ME"
var OUTDATED = "OUTDATED"
var SYNC = "SYNC"

func checkError(err error) {
	if err != nil {
		println(err)
	}
}

func TranslateMessage(message commons.Message, port string) {

	if message.Port == port {
		return
	}
	connection := CreateConnection()
	fmt.Println(message, port)

	channel := CreateChannel(connection, message.Port)
	_prt := message.Port
	message.Port = port

	if message.RequestType == UPDATE_ME {
		var datetime time.Time
		// if the time is not specified we use this as a default
		datetime, err := time.Parse("2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05.000Z")
		checkError(err)
		if len(message.Operations) > 0 {
			err := json.Unmarshal(message.Operations[0].Data, &datetime)
			fmt.Println("DATE TIME ", datetime)
			checkError(err)
		}

		operations := commons.GetOperations(datetime)
		message.Operations = operations
		message.RequestType = YOUR_UPDATE
		msg, err := json.Marshal(message)
		checkError(err)
		PublishMessage(channel, msg, _prt)
		return
	}

	if message.RequestType == YOUR_UPDATE && len(message.Operations) != 0 {
		latestOps := commons.GetLatestOperation()
		fmt.Println(latestOps)
		if latestOps.Sequence.After(message.Operations[0].Sequence) { // if there is latest operation after the request
			message.RequestType = OUTDATED
			msg, err := json.Marshal(message)
			checkError(err)
			PublishMessage(channel, msg, _prt)
			return
		}

		for _, operation := range message.Operations {
			if operation.DataType == PROJECT {
				projectHandler(operation)
			}
			if operation.DataType == TASK {
				taskHandler(operation)
			}
		}
		return
	}

	if message.RequestType == OUTDATED {
		message.RequestType = UPDATE_ME
		msg, err := json.Marshal(message)
		checkError(err)
		PublishMessage(channel, msg, _prt)
		return
	}

	if message.RequestType == SYNC {
		operation := message.Operations[0]
		if operation.DataType == "PROJECT" {
			projectHandler(operation)
		}
		if operation.DataType == "TASK" {
			taskHandler(operation)
		}
	}
}

func PropagateUpdate(msg []byte, port string) {
	connection := CreateConnection()
	for _, server := range Servers {
		if port != server { // avoids sending to itself
			channel := CreateChannel(connection, server)
			PublishMessage(channel, msg, port)
			_ = channel.Close()
		}
	}
}

func projectHandler(operation commons.Operation) {
	var project models.Project
	err := json.Unmarshal(operation.Data, &project)
	checkError(err)
	println("hmm")
	fmt.Println(project)
	if operation.OpType == CREATE {
		models.CreateProject(project)
	}

	if operation.OpType == UPDATE {
		models.UpdateProject(string(project.ID), project)
	}

	if operation.OpType == DELETE {
		models.DeleteProject(string(project.ID))
	}
}

func taskHandler(operation commons.Operation) {
	var task models.Task
	err := json.Unmarshal(operation.Data, &task)
	checkError(err)
	if operation.OpType == CREATE {
		models.CreateTask(task)
	}

	if operation.OpType == UPDATE {
		models.UpdateTask(string(task.ID), task)
	}

	if operation.OpType == DELETE {
		models.DeleteTask(string(task.ID))
	}
}

func UpdateMe(channel *amqp.Channel, port string) {
	if len(Servers) == 0 || port == Servers[0] {
		return
	}
	message := commons.Message{
		RequestType: "UPDATE-ME",
		Port:        port,
		Operations:  nil,
	}
	message.Operations = append(message.Operations, commons.GetLatestOperation())
	msg, err := json.Marshal(message)
	checkError(err)
	PublishMessage(channel, msg, Servers[0])
}
