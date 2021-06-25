package rabbitmq

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	delivery <-chan amqp.Delivery
}

// Connects and consumes from rabbit queue
func New(queueAddress string, queueName string) (*Consumer, error) {

	c := &Consumer{
		conn:    nil,
		channel: nil,
	}

	var err error

	log.Printf("Trying to connect to the RabbitMQ queue %s, at address %s", queueName, queueAddress)
	c.conn, err = amqp.Dial(queueAddress)
	logger.FailOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	c.channel, err = c.conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")
	//defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("WARNING: Issues found declaring queue: %s", err)
	}

	c.delivery, err = c.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	logger.FailOnError(err, "Failed on reading")

	return c, err
}