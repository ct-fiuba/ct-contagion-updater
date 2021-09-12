package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/controllers"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/contagions"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/concurrency"
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

	visitsCollection, err := visits.New(db)
	logger.FailOnError(err, "Failed to create/get visits collection")
	visits, err := visitsCollection.All()
	fmt.Printf("### VISITS: \n%+v\n", visits)

	infectedManager := controllers.NewInfectedManager(db)

	forever := make(chan bool)

	go func() {
		codesBySpace := concurrency.NewSafeStringListMap()
		for d := range consumer.DeliveryChan {
			fmt.Printf(">>> Consumed a message: %s\n", d.Body)

			var contagion contagions.Contagion
			err := json.Unmarshal(d.Body, &contagion)
			if err != nil {
				log.Printf("[ERROR] Unmarshall failure %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("syntax error at byte offset %d", e.Offset)
				}
			} else {
				fmt.Printf(">>>>>> To struct: %+v\n", contagion)
				codesBySpace.Add(contagion.SpaceId, contagion.UserGeneratedCode)
				infectedManager.ProcessBatch(codesBySpace)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
