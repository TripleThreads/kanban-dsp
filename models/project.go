package models

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "kanban-distributed-system/commons"
	. "kanban-distributed-system/config"
	"kanban-distributed-system/utility"
	"strconv"
	"time"
)

type Project struct {
	gorm.Model
	Title string `json:"title"`
	Tasks []Task
}

func checkError(err error) {
	if err != nil {
		println(err)
	}
}

func GetProjects() []Project {
	db := ConnectToMysql()
	var projects []Project
	db.Find(&projects)
	_ = db.Close()
	return projects
}

func GetProject(id string) Project {
	var project Project
	var tasks []Task

	db := ConnectToMysql()
	println(id)
	db.Where("id = ?", id).Find(&project)
	println(id)
	db.Where("project_id = ?", id).Find(&tasks)
	project.Tasks = tasks
	fmt.Println(tasks)
	_ = db.Close()
	return project
}

func CreateProject(project Project, timestamp time.Time) []byte {
	db := ConnectToMysql()
	db.Create(&project)
	body, err := json.Marshal(project)
	checkError(err)
	// LOG CURRENT OPERATION
	msg := LogOperation(body, "CREATE", "PROJECT", timestamp)
	_ = db.Close()
	return msg
}

func UpdateProject(id string, project Project, timestamp time.Time) []byte {
	var pr Project

	db := ConnectToMysql()
	db.Where("id = ?", id).Find(&pr)
	db.Model(pr).Updates(project)

	ID, _ := strconv.Atoi(id)
	project.ID = uint(ID)

	body, err := json.Marshal(project)
	checkError(err)
	// LOG CURRENT OPERATION
	msg := LogOperation(body, "UPDATE", "PROJECT", timestamp)
	_ = db.Close()
	return msg
}

func DeleteProject(id string, timestamp time.Time) []byte {
	var project Project

	db := ConnectToMysql()
	db.Where("id = ?", id).Find(&project)
	db.Delete(project)

	ID, _ := strconv.Atoi(id)
	project.ID = uint(ID)

	println("id number ", ID)
	body, err := json.Marshal(project)
	checkError(err)
	// LOG CURRENT OPERATION
	msg := LogOperation(body, "DELETE", "PROJECT", timestamp)
	_ = db.Close()
	return msg
}

func LogOperation(body []byte, OpType string, DataType string, timestamp time.Time) []byte {
	var operations []Operation
	operations = append(operations, Operation{Data: body, OpType: OpType, DataType: DataType, Sequence: timestamp})
	message := Message{RequestType: "SYNC", Operations: operations, Port: utility.PORT}
	msg, _ := json.Marshal(message)
	CreateOperation(operations[0])
	return msg
}
