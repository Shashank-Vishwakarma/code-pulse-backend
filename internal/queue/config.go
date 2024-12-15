package queue

import (
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Ch *amqp.Channel
var Q *amqp.Queue

func InitializeRabbitMQ() error {
	// create a connection
	conn, err := amqp.Dial(config.Config.RABBITMQ_URL)
	if err != nil {
		return err
	}

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	Ch = ch

	// create a queue
	q, err := ch.QueueDeclare(config.Config.QUEUE_NAME, false, false, false, false, nil)
	if err != nil {
		return err
	}
	Q = &q

	return nil
}
