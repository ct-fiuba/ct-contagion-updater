package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/contagions"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	// rd_api "github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	rd_impl "github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/impl"
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

func (self *InfectedManager) Process() error {
	rulesCollection, err := rules.New(self.db)
	logger.FailOnError(err, "Failed to create/get rules collection")
	rules, err := rulesCollection.All()

	visitsCollection, err := visits.New(self.db)
	logger.FailOnError(err, "Failed to create/get visits collection")

	spacesCollection, err := spaces.New(self.db)
	logger.FailOnError(err, "Failed to create/get spaces collection")

	ruleChain := rd_impl.NewSimpleRuleChain(visitsCollection, spacesCollection)
	for i, rule := range rules {
		fmt.Printf("RULE #%d --> %+v\n", i, rule)
		ruleChain.AddFilter(fmt.Sprint(i), rule)
	}

	for d := range self.queue.DeliveryChan {
		var contagion contagions.Contagion
		err := json.Unmarshal(d.Body, &contagion)
		if err != nil {
			logger.FailOnError(err, "Failed unmarshalling contagion")
		} else {
			fmt.Printf(">>> Consumed a message: %s\n", d.Body)
			fmt.Printf(">>>>>> To struct: %+v\n", contagion)

			ruleChain.Process(contagion)
		}
	}
	return nil
}
