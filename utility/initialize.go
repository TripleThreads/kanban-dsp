package utility

import (
	"bytes"
	"encoding/json"
	. "kanban-distributed-system/commons"
	"net/http"
)

var PORT string

func RegisterServer(port string) string {

	PORT = port
	msg := Message{Port: PORT, RequestType: "REG"}

	encoded, _ := json.Marshal(msg)

	http.Post("http://127.0.0.1:9865", "application/json", bytes.NewBuffer(encoded))

	println("Running on :", PORT)

	return PORT
}
