package services

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "search_queue"
)

func PublishHotel(hotel interface{}) error {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error opening RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error declaring RabbitMQ queue: %v", err)
	}

	body, err := json.Marshal(hotel)
	if err != nil {
		return fmt.Errorf("error marshalling hotel data: %v", err)
	}

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	return err
}
