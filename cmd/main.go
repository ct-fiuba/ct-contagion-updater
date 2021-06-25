package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/rabbitmq"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"

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

	consumer, err := rabbitmq.New(queueAddress, queueName)
	logger.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range consumer.delivery {
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
