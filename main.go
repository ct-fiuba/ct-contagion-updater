package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	queueAddress := os.Getenv("QUEUE_ADDRESS")
	queueName := os.Getenv("QUEUE_NAME")
	log.Printf("Trying to connect to the RabbitMQ queue %s, at address %s", queueName, queueAddress)

	conn, err := amqp.Dial(queueAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// _, err = ch.QueueDeclare(
	// 	queueName, // name
	// 	false,     // durable
	// 	false,     // delete when unused
	// 	false,     // exclusive
	// 	false,     // no-wait
	// 	nil,       // arguments
	// )
	// if err != nil {
	// 	log.Printf("WARNING: Issues found declaring queue: %s", err)
	// }

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(">>> Consumed a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
