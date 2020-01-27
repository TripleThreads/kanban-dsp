package message_queue

import (
	"encoding/json"
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

func checkError(err error) {
	if err != nil {
		println(err)
	}
}

func TranslateMessage(message commons.Message, port string) {

	channel := CreateChannel(message.Port)

	if message.RequestType == UPDATE_ME {
		var datetime time.Time
		err := json.Unmarshal(message.Operations[0].Data, &datetime)

		checkError(err)

		operations := commons.GetOperations(datetime)

		message.Operations = operations
		message.Port = port
		message.RequestType = YOUR_UPDATE
		msg, err := json.Marshal(message)
		checkError(err)
		PublishMessage(msg, channel)
	}

	if message.RequestType == YOUR_UPDATE {
		latestOps := commons.GetLatestOperation()

		if latestOps.Sequence.Before(message.Operations[0].Sequence) { // if there is latest operation after the request
			message.Port = port
			message.RequestType = OUTDATED
			msg, err := json.Marshal(message)
			checkError(err)
			PublishMessage(msg, channel)
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
	}

	if message.RequestType == OUTDATED {
		message.RequestType = UPDATE_ME
		msg, err := json.Marshal(message)
		checkError(err)
		PublishMessage(msg, channel)
	}
}

func PropagateUpdate(msg []byte, port string) {
	for _, server := range Servers {
		if port != server { // avoids sending to itself
			channel := CreateChannel(server)
			PublishMessage(msg, channel)
		}
	}
}

func projectHandler(operation commons.Operation) {
	var project models.Project
	err := json.Unmarshal(operation.Data, &project)
	checkError(err)
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
