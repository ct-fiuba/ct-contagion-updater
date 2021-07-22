package impl

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
)

type SimpleRuleChain struct {
	filters       []*api.RuleChecker
	resultHandler *api.ResultHandler
}

func NewSimpleRuleChain(rh *api.ResultHandler) api.RuleChain {
	self := new(SimpleRuleChain)
	self.resultHandler = rh

	return self
}

func (self *SimpleRuleChain) AddFilter(id string, filter api.RuleChecker) bool {
	return true
}

func (self *SimpleRuleChain) RemoveFilter(id string) bool {
	return true
}

func (self *SimpleRuleChain) Process(visit visits.Visit) {
	return
}
