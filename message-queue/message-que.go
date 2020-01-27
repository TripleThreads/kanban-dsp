package message_queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	. "kanban-distributed-system/commons"
)

var Servers []string

func CreateChannel(bindAddress string) *amqp.Channel {
	println("creating channel ", bindAddress)
	url := "amqp://guest:guest@localhost:5672"

	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("Couldn't create channel" + err.Error())
	}

	// We create an exchange that will bind to the queue to send and receive messages
	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	// We create a queue named Test
	_, err = channel.QueueDeclare(bindAddress, true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	err = channel.QueueBind(bindAddress, "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	return channel
}

func PublishMessage(msg []byte, channel *amqp.Channel) {
	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct
	message := amqp.Publishing{
		Body: msg,
	}

	// We publish the message to the exchange we created earlier
	err := channel.Publish("events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}
}

func ConsumeMessage(channel *amqp.Channel, port string) {
	println("port", port)
	// We consume data from the queue named Test using the channel we created in go.
	msgs, err := channel.Consume(port, "", false, false, false,
		false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	// We loop through the messages in the queue and print them in the console.
	// The msgs will be a go channel, not an amqp channel
	var message Message
	for msg := range msgs {
		_ = json.Unmarshal(msg.Body, &message)

		if message.RequestType == "SERVERS" {
			_ = json.Unmarshal(message.Operations[0].Data, &Servers)
		} else {
			TranslateMessage(message, port)
		}
		_ = msg.Ack(false)
	}
}
