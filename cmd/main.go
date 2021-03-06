package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/controllers"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/contagions"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/concurrency"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/rabbitmq"

	cron "github.com/robfig/cron/v3"
)

const DEFAULT_SCHEDULE_PERIOD = "@every 20s"

func main() {
	queueAddress := os.Getenv("QUEUE_ADDRESS")
	queueName := os.Getenv("QUEUE_NAME")
	dbUri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")

	schedulePeriod := os.Getenv("TRIGGER_SCHEDULE_PERIOD")
	if schedulePeriod == "" {
		schedulePeriod = DEFAULT_SCHEDULE_PERIOD
	}

	db, err := mongodb.New(dbUri, dbName)
	logger.FailOnError(err, "Failed to connect to the DB")
	defer db.Shutdown()

	consumer, err := rabbitmq.New(queueAddress, queueName)
	logger.FailOnError(err, "Failed to register a consumer")
	defer consumer.Shutdown()

	infectedManager := controllers.NewInfectedManager(db)
	codesBySpace := concurrency.NewSafeStringListMap()

	forever := make(chan bool)

	c := cron.New()
	c.AddFunc(schedulePeriod, func() {
		log.Printf("[MAIN] Starting batch processing\n")
		infectedManager.ProcessBatch(codesBySpace)
	})
	c.Start()

	go func() {
		for d := range consumer.DeliveryChan {
			log.Printf("[MAIN] Consumed a message: %s\n", d.Body)

			var contagion contagions.Contagion
			err := json.Unmarshal(d.Body, &contagion)
			if err != nil {
				log.Printf("[ERROR] Unmarshall failure %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("syntax error at byte offset %d", e.Offset)
				}
			} else {
				log.Printf("[MAIN] >>> To struct: %+v\n", contagion)
				codesBySpace.Add(contagion.SpaceId, contagion.UserGeneratedCode)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
