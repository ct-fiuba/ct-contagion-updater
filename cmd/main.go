package main

import (
	"log"
	"os"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/rabbitmq"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
)



func main() {
	queueAddress := os.Getenv("QUEUE_ADDRESS")
	queueName := os.Getenv("QUEUE_NAME")
	dbUri := os.Getenv("MONGODB_URI")

	db, err := mongodb.New(dbUri, "contact-tracing-db")
	logger.FailOnError(err, "Failed to register a consumer")
	defer db.Shutdown()

	consumer, err := rabbitmq.New(queueAddress, queueName)
	logger.FailOnError(err, "Failed to register a consumer")
	defer consumer.Shutdown()

	forever := make(chan bool)

	go func() {
		rules, err := rules.New(db)
		logger.FailOnError(err, "Failed create rules collection")
		rules.All()
		for d := range consumer.Delivery {
			log.Printf(">>> Consumed a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
