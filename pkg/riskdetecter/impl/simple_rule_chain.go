package impl

import (
	"fmt"
	"time"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/contagions"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FilterSlot struct {
	id      string
	checker api.RuleChecker
}

type SimpleRuleChain struct {
	filters       []FilterSlot
	resultHandler *SimpleResultHandler
	visitsColl    *visits.VisitsCollection
	spacesColl    *spaces.SpacesCollection
}

func NewSimpleRuleChain(visitsColl *visits.VisitsCollection, spacesColl *spaces.SpacesCollection) api.RuleChain {
	self := new(SimpleRuleChain)
	self.resultHandler = NewSimpleResultHandler()
	self.visitsColl = visitsColl
	self.spacesColl = spacesColl

	return self
}

func (self *SimpleRuleChain) AddFilter(id string, rule rules.Rule) bool {
	filter := NewSimpleRuleChecker(rule)
	if len(self.filters) > 0 {
		lastFilter := self.filters[len(self.filters)-1]
		lastFilter.checker.SetOutput(filter.AsOutputConnector())
	}

	fs := FilterSlot{id: id, checker: filter.AsRuleChecker()}
	self.filters = append(self.filters, fs)

	rc := self.resultHandler.AsResultConnector()
	filter.SetResultExit(rc)
	return true
}

func (self *SimpleRuleChain) RemoveFilter(id string) bool {
	return false
}

func (self *SimpleRuleChain) Process(contagion contagions.Contagion) error {
	someTime := time.Now()
	someDuration, _ := time.ParseDuration("10m")
	someDuration2, _ := time.ParseDuration("20m")

	allVisits, err := self.visitsColl.All()
	logger.FailOnError(err, "Failed retrieving visits!")
	fmt.Printf("VISITS = %v \n", allVisits)

	space, err := self.spacesColl.Find(contagion.SpaceId)
	logger.FailOnError(err, "Failed retrieving space!")
	fmt.Printf("SPACE = %+v \n", space)

	infectedVisit := visits.Visit{
		ScanCode:          primitive.NewObjectID(),
		UserGeneratedCode: "Hola",
		EntranceTimestamp: primitive.NewDateTimeFromTime(someTime),
		Vaccinated:        0,
		CovidRecovered:    false,
	}

	relatedVisits := []visits.Visit{
		{
			ScanCode:          primitive.NewObjectID(),
			UserGeneratedCode: "Hola11",
			EntranceTimestamp: primitive.NewDateTimeFromTime(someTime.Add(-someDuration)),
			Vaccinated:        0,
			CovidRecovered:    false,
		},
		{
			ScanCode:          primitive.NewObjectID(),
			UserGeneratedCode: "Hola22",
			EntranceTimestamp: primitive.NewDateTimeFromTime(someTime.Add(-someDuration2)),
			Vaccinated:        0,
			CovidRecovered:    false,
		},
	}

	for _, v := range relatedVisits {
		initialFilter := self.filters[0].checker
		fmt.Printf("AAAAAAAAAA, %+v \n", initialFilter)
		initialFilter.Execute(infectedVisit, v)
	}

	return nil
}
