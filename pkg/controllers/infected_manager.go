package controllers

import (
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	// rd_api "github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	rd_impl "github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/impl"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/concurrency"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"
)

type InfectedManager struct {
	db *mongodb.DB
}

func NewInfectedManager(db *mongodb.DB) *InfectedManager {
	return &InfectedManager{
		db: db,
	}
}

func (self *InfectedManager) ProcessBatch(codesBySpace *concurrency.SafeStringListMap) error {
	if codesBySpace.Count() == 0 {
		fmt.Printf("[DEBUG] Empty batch. Skipping processing\n")
		return nil
	}

	rulesCollection, err := rules.New(self.db)
	logger.FailOnError(err, "Failed to create/get rules collection")
	visitsCollection, err := visits.New(self.db)
	logger.FailOnError(err, "Failed to create/get visits collection")
	spacesCollection, err := spaces.New(self.db)
	logger.FailOnError(err, "Failed to create/get spaces collection")
	compromisedCodesCollection, err := compromisedCodes.New(self.db)
	logger.FailOnError(err, "Failed to create/get spaces collection")

	rules, err := rulesCollection.All()
	fmt.Printf("[DEBUG] Rules list: %+v\n", rules)
	ruleChain := rd_impl.NewSimpleRuleChain(compromisedCodesCollection)
	for i, rule := range rules {
		ruleChain.AddFilter(fmt.Sprint(i), rule)
	}

	codesBySpaceBatch := codesBySpace.Clear()

	for spaceId, codes := range codesBySpaceBatch {
		fmt.Printf("[DEBUG] Processing codes %+v for space %s\n", codes, spaceId)

		visits, err := visitsCollection.FindInSpace(spaceId)
		if err != nil {
			fmt.Printf("[ERROR] Failure when retrieving visits for space %s \n", spaceId)
			continue
		}
		fmt.Printf("[DEBUG] Visits collected = %d\n", len(visits))

		for _, contagionCode := range codes {
			contagionVisit, err := visitsCollection.FindByGeneratedCode(contagionCode)
			if err != nil {
				fmt.Printf("[ERROR] Failure when retrieving visit info from contagion with code %s \n", contagionCode)
				continue
			}

			space, err := spacesCollection.Find(spaceId)
			if err != nil {
				fmt.Printf("[ERROR] Failure when retrieving space %s \n", spaceId)
				continue
			}

			for _, v := range visits {
				// TODO: Try to avoid processing already infected visits
				if v.UserGeneratedCode == contagionVisit.UserGeneratedCode {
					continue
				}

				fmt.Printf("[TRACE] Adding to chain: \n\tProcessed Visit: %s\n\tContagion Visit: %s\n\tSpace: %s\n", v.UserGeneratedCode, contagionVisit.UserGeneratedCode, spaceId)
				err := ruleChain.Process(&v, contagionVisit, space)
				if err != nil {
					fmt.Printf("[ERROR] Failure while procesing visits\n\tProcessed Visit: %+v\n\tContagion Visit: %+v\n\tSpace: %+v\n", v, contagionVisit, space)
					continue
				}
			}
		}
	}

	return nil
}
