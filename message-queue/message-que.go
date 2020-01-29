package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	. "kanban-distributed-system/commons"
)

var Servers []string
var CONNECTION *amqp.Connection

func CreateConnection() *amqp.Connection {
	url := "amqp://guest:guest@localhost:5672"
	var err error
	CONNECTION, err = amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	return CONNECTION
}

func CreateChannel(connection *amqp.Connection, bindAddress string) *amqp.Channel {

	println("creating channel ", bindAddress)
	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("Couldn't create channel" + err.Error())
	}

	// We create an exchange that will bind to the queue to send and receive messages
	err = channel.ExchangeDeclare(bindAddress, "direct", true, false, false, false, nil)

	// We create a queue named Test
	_, err = channel.QueueDeclare(bindAddress, true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	err = channel.QueueBind(bindAddress, "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	return channel
}

func PublishMessage(channel *amqp.Channel, msg []byte, port string) {

	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct
	message := amqp.Publishing{
		Body: msg,
	}

	// We publish the message to the exchange we created earlier
	err := channel.Publish("events", port, false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}
}

func ConsumeMessage(channel *amqp.Channel, port string) {

	// We consume data from the queue named Test using the channel we created in go.
	msgs, err := channel.Consume(port, "", true, false, false,
		false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	var message Message
	for msg := range msgs {
		err = json.Unmarshal(msg.Body, &message)
		checkError(err)
		if message.RequestType == "SERVERS" {
			_ = json.Unmarshal(message.Operations[0].Data, &Servers)
		} else {
			TranslateMessage(message, port)
		}
	}
}
