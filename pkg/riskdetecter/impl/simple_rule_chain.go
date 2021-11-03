package impl

import (
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	timeutils "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/time"
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

func (rulechain *SimpleRuleChain) Process(compromised, infected *visits.Visit, s *spaces.Space) error {
	compromisedEntranceTime, compromisedExitTime := GetVisitInterval(compromised, s.EstimatedVisitDuration)
	infectedEntranceTime, infectedExitTime := GetVisitInterval(infected, s.EstimatedVisitDuration)
	fmt.Printf("[Rule Chain] Compromised time interval = [ %s ; %s ] \n", compromisedEntranceTime.String(), compromisedExitTime.String())
	fmt.Printf("[Rule Chain] Infected time interval = [ %s ; %s ] \n", infectedEntranceTime.String(), infectedExitTime.String())

	visitsOverlapping := timeutils.IntervalsOverlap(compromisedEntranceTime, compromisedExitTime, infectedEntranceTime, infectedExitTime)

	if !visitsOverlapping {
		fmt.Printf("[Rule Chain] Visits not overlapping. Skipping checks.\n")
		return nil
	}

	initialFilter := rulechain.filters[0].checker
	err := initialFilter.Process(compromised, infected, s)
	return err
}
