package main

import (
	"log"
	"os"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/controllers"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/rabbitmq"
)

func main() {
	queueAddress := os.Getenv("QUEUE_ADDRESS")
	queueName := os.Getenv("QUEUE_NAME")
	dbUri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")

	db, err := mongodb.New(dbUri, dbName)
	logger.FailOnError(err, "Failed to connect to the DB")
	defer db.Shutdown()

	consumer, err := rabbitmq.New(queueAddress, queueName)
	logger.FailOnError(err, "Failed to register a consumer")
	defer consumer.Shutdown()

	infectedManager := controllers.NewInfectedManager(consumer, db)

	forever := make(chan bool)

	go func() {
		infectedManager.Process()
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
