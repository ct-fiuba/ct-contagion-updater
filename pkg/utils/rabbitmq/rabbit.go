package rabbitmq

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	Conn     *amqp.Connection
	Channel  *amqp.Channel
	Delivery <-chan amqp.Delivery
}

// Connects and consumes from rabbit queue
func New(queueAddress string, queueName string) (*Consumer, error) {
	c := &Consumer{
		Conn:     nil,
		Channel:  nil,
		Delivery: nil,
	}

	var err error

	log.Printf("Trying to Connect to the RabbitMQ queue %s, at address %s", queueName, queueAddress)
	c.Conn, err = amqp.Dial(queueAddress)
	logger.FailOnError(err, "Failed to Connect to RabbitMQ")
	//defer Conn.Close()

	c.Channel, err = c.Conn.Channel()
	logger.FailOnError(err, "Failed to open a Channel")
	//defer ch.Close()

	_, err = c.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	logger.FailOnError(err, "WARNING: Issues found declaring queue")

	c.Delivery, err = c.Channel.Consume(
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

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel

	if err := c.Conn.Close(); err != nil {
		logger.FailOnError(err, "Failed to close connection")
		return err
	}

	if err := c.Channel.Close(); err != nil {
		logger.FailOnError(err, "Failed to close channel")
		return err
	}

	defer log.Printf("AMQP shutdown OK")
	return nil
}
