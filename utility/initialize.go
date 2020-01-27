package utility

import (
	"bytes"
	"encoding/json"
	"github.com/streadway/amqp"
	. "kanban-distributed-system/commons"
	. "kanban-distributed-system/message-queue"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var PORT string
var ReadChannel *amqp.Channel
var WriteChannel *amqp.Channel

func RegisterServer() string {
	rand.Seed(time.Now().UnixNano())

	port := 10000 + rand.Intn(1000)
	PORT = ":" + strconv.FormatInt(int64(port), 10)

	msg := Message{Port: PORT, RequestType: "REG"}
	WriteChannel = CreateChannel(PORT)
	encoded, _ := json.Marshal(msg)

	go PublishMessage(encoded, WriteChannel)

	ReadChannel = CreateChannel(PORT)

	go ConsumeMessage(WriteChannel, PORT)

	http.Post("http://127.0.0.1:9865", "application/json", bytes.NewBuffer(encoded))

	println("Running on :", PORT)

	return PORT
}
