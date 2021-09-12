package impl

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
)

type FilterSlot struct {
	id      string
	checker api.RuleChecker
}

type SimpleRuleChain struct {
	filters       []FilterSlot
	resultHandler *SimpleResultHandler
}

func NewSimpleRuleChain(compromisedCodesCollection *compromisedCodes.CompromisedCodesCollection) api.RuleChain {
	rulechain := new(SimpleRuleChain)
	rulechain.resultHandler = NewSimpleResultHandler(compromisedCodesCollection)

	return rulechain
}

func (rulechain *SimpleRuleChain) AddFilter(id string, rule rules.Rule) bool {
	filter := NewSimpleRuleChecker(rule)
	if len(rulechain.filters) > 0 {
		lastFilter := rulechain.filters[len(rulechain.filters)-1]
		lastFilter.checker.SetNext(filter)
	}

	fs := FilterSlot{id: id, checker: filter}
	rulechain.filters = append(rulechain.filters, fs)

	rc := rulechain.resultHandler.AsResultConnector()
	filter.SetResultExit(rc)
	return true
}

func (rulechain *SimpleRuleChain) RemoveFilter(id string) bool {
	return false // TODO
}

func (rulechain *SimpleRuleChain) Process(v1, v2 *visits.Visit, s *spaces.Space) error {
	// someTime := time.Now()
	// someDuration, _ := time.ParseDuration("10m")
	// someDuration2, _ := time.ParseDuration("20m")

	// infectedVisit := visits.Visit{
	// 	ScanCode:          primitive.NewObjectID(),
	// 	UserGeneratedCode: "Hola",
	// 	EntranceTimestamp: primitive.NewDateTimeFromTime(someTime),
	// 	Vaccinated:        0,
	// 	CovidRecovered:    false,
	// }

	// relatedVisits := []visits.Visit{
	// 	{
	// 		ScanCode:          primitive.NewObjectID(),
	// 		UserGeneratedCode: "Hola11",
	// 		EntranceTimestamp: primitive.NewDateTimeFromTime(someTime.Add(-someDuration)),
	// 		Vaccinated:        0,
	// 		CovidRecovered:    false,
	// 	},
	// 	{
	// 		ScanCode:          primitive.NewObjectID(),
	// 		UserGeneratedCode: "Hola22",
	// 		EntranceTimestamp: primitive.NewDateTimeFromTime(someTime.Add(-someDuration2)),
	// 		Vaccinated:        0,
	// 		CovidRecovered:    false,
	// 	},
	// }
	// for _, v := range relatedVisits {
	// 	initialFilter := rulechain.filters[0].checker
	// 	fmt.Printf("AAAAAAAAAA, %+v \n", initialFilter)
	// 	initialFilter.Execute(infectedVisit, v)
	// }

	initialFilter := rulechain.filters[0].checker
	err := initialFilter.Process(v1, v2, s)
	return err
}
