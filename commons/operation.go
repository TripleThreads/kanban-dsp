package commons

import (
	"fmt"
	"kanban-distributed-system/config"
	"time"
)

type Operation struct {
	id       int       `gorm:"AUTO_INCREMENT;column:id;primary_key"`
	Sequence time.Time `gorm:"unique"`
	OpType   string    // UPDATE, CREATE, DELETE
	DataType string    // is it task or project
	Data     []byte    `gorm:"unique"` // marshalled data
}

func GetOperations(datetime time.Time) []Operation {
	db := config.ConnectToMysql()
	var operations []Operation
	db.Raw("SELECT * FROM operations WHERE Date(sequence) > Date(?) order by sequence asc", datetime).Scan(&operations)
	fmt.Println("latest ", operations)
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
	db.Order("sequence desc").First(&operation)
	return operation
}
