package impl

import (
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
)

type SimpleRuleChecker struct {
	rule       rules.Rule
	next       api.OutputConnector
	resultExit api.ResultConnector
}

func NewSimpleRuleChecker(r rules.Rule) api.RuleChecker {
	self := new(SimpleRuleChecker)
	self.rule = r

	return self
}

func (self *SimpleRuleChecker) SetInput(ic api.InputConnector) {
	// NOOP
}

func (self *SimpleRuleChecker) SetOutput(oc api.OutputConnector) {
	self.next = oc
}

func (self *SimpleRuleChecker) SetResultExit(rc api.ResultConnector) {
	self.resultExit = rc
}

func (self *SimpleRuleChecker) Execute(visit visits.Visit) {
	return
}

func (self *SimpleRuleChecker) Push(v visits.Visit) bool {
	self.Execute(v)
	return true
}

// -- CASTERS
func (self *SimpleRuleChecker) asOutputConnector() api.OutputConnector {
	return self
}
