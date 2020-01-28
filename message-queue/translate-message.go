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
var UPDATED = "UPDATED"

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

		if len(message.Operations) > 0 {
			datetime = message.Operations[0].Sequence
		}

		fmt.Println("DATE TIME ", datetime)

		operations := commons.GetOperations(datetime)
		message.RequestType = YOUR_UPDATE
		fmt.Println("UPDATES ", operations)
		if len(operations) == 0 {
			println("YOU ARE UPTO DATE... ")
			message.RequestType = UPDATED
		}
		message.Operations = operations
		msg, err := json.Marshal(message)
		checkError(err)
		PublishMessage(channel, msg, _prt)
		return
	}

	if message.RequestType == YOUR_UPDATE && len(message.Operations) != 0 {
		latestOp := commons.GetLatestOperation()
		fmt.Println(latestOp)

		// if there is latest operation after the request
		if latestOp.Sequence.Equal(message.Operations[0].Sequence) || latestOp.Sequence.After(message.Operations[0].Sequence) {
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
		latestOp := commons.GetLatestOperation()
		operation := message.Operations[0]
		if latestOp.Sequence.Equal(operation.Sequence) || latestOp.Sequence.After(operation.Sequence) {
			return
		}
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
		models.CreateProject(project, operation.Sequence)
	}

	if operation.OpType == UPDATE {
		models.UpdateProject(string(project.ID), project, operation.Sequence)
	}

	if operation.OpType == DELETE {
		models.DeleteProject(string(project.ID), operation.Sequence)
	}
}

func taskHandler(operation commons.Operation) {
	var task models.Task
	err := json.Unmarshal(operation.Data, &task)
	checkError(err)
	if operation.OpType == CREATE {
		models.CreateTask(task, operation.Sequence)
	}

	if operation.OpType == UPDATE {
		models.UpdateTask(string(task.ID), task, operation.Sequence)
	}

	if operation.OpType == DELETE {
		models.DeleteTask(string(task.ID), operation.Sequence)
	}
}

func UpdateMe(channel *amqp.Channel, port string) {
	fmt.Println(port, Servers)
	latest := commons.GetLatestOperation()
	if len(Servers) == 0 || port == Servers[0] {
		return
	}
	message := commons.Message{
		RequestType: UPDATE_ME,
		Port:        port,
		Operations:  nil,
	}
	message.Operations = append(message.Operations, latest)
	fmt.Println("LATEST", message.Operations[0].Sequence)
	msg, err := json.Marshal(message)
	checkError(err)
	PublishMessage(channel, msg, Servers[0])
}
