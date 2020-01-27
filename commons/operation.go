package commons

import (
	"fmt"
	"kanban-distributed-system/config"
	"time"
)

type Operation struct {
	Sequence time.Time
	OpType   string // UPDATE, CREATE, DELETE
	DataType string // is it task or project
	Data     []byte // marshalled data
}

func GetOperations(datetime time.Time) []Operation {
	db := config.ConnectToMysql()
	var operations []Operation
	db.Where("sequence > ?", datetime).Find(&operations)
	_ = db.Close()
	return operations
}

func CreateOperation(operation Operation) {
	db := config.ConnectToMysql()
	db.Create(operation)
	fmt.Println("operation created successfully")
}

func GetLatestOperation() Operation {
	db := config.ConnectToMysql()
	var operation Operation
	db.Last(&operation)
	return operation
}
