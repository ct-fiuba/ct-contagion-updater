package controllers

import (
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/rabbitmq"
)

type InfectedManager struct {
	queue *rabbitmq.Consumer
	db    *mongodb.DB
}

func NewInfectedManager(queue *rabbitmq.Consumer, db *mongodb.DB) *InfectedManager {
	return &InfectedManager{
		queue: queue,
		db:    db,
	}
}

func (im *InfectedManager) Process() error {
	rules, err := rules.New(im.db)
	visits, err := visits.New(im.db)
	logger.FailOnError(err, "Failed to create/get rules collection")
	rules.All()
	visits.All()
	for d := range im.queue.Delivery {
		log.Printf(">>> Consumed a message: %s", d.Body)
	}
	return nil
}
