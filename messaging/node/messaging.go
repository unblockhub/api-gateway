package node

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var (
	logger           = log.New(os.Stdout, "[GATEWAY][MESSAGING][INTERNAL]", 0)
	rabbit           *amqp.Connection
	rabbitConnection *amqp.Channel
	queues = [1]string{"API_DISCONNECT"}
)

func Init() {
	rabbitClient, err := amqp.Dial(os.Getenv("RABBIT_URI"))
	if err != nil {
		panic(fmt.Sprintf("Rabbit could not be reached: %s", err))
	}
	rabbit = rabbitClient
	rabbitConnection, _ = rabbit.Channel()

	// Create queues
	for i := 0; i < len(queues); i++ {
		queueName := queues[i]
		_, _ = rabbitConnection.QueueDeclare(queueName, false, false, true, true, nil)
	}
}

func Subscribe(channel string, callback func([]byte)) {
	msgs, err := rabbitConnection.Consume(
		channel,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Printf("Failed to subscribe to channel: %s", err)
	}

	go func() {
		for message := range msgs {
			body := message.Body

			go callback(body)
		}
	}()
}

func Publish(content interface{}, channel string) {
	encodedContent, _ := json.Marshal(content)
	err := rabbitConnection.Publish(
		"",
		channel,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encodedContent,
		},
	)
	if err != nil {
		logger.Printf("Failed to publish message: %s", err)
		return
	}
}
