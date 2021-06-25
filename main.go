package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rule struct {
	Index         int    `bson:"index"`
	ContagionRisk string `bson:"contagionRisk"`
	DurationCmp   string `bson:"durationCmp,omitempty"`
	DurationValue int    `bson:"durationValue,omitempty"`
	M2Cmp         string `bson:"m2Cmp,omitempty"`
	M2Value       int    `bson:"m2Value,omitempty"`
	SpaceValue    string `bson:"spaceValue,omitempty"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	queueAddress := os.Getenv("QUEUE_ADDRESS")
	queueName := os.Getenv("QUEUE_NAME")
	dbUri := os.Getenv("MONGODB_URI")

	ctx, cancel := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	defer cancel()

	log.Printf("Connected to the DB!")

	database := client.Database("contact-tracing-db")
	rulesCollection := database.Collection("rules")

	log.Printf("Trying to connect to the RabbitMQ queue %s, at address %s", queueName, queueAddress)
	conn, err := amqp.Dial(queueAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

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
			var rules []Rule
			cursor, err := rulesCollection.Find(ctx, bson.D{})
			if err != nil {
				panic(err)
			}
			if err = cursor.All(ctx, &rules); err != nil {
				panic(err)
			}
			fmt.Println(rules)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
