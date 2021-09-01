package impl

import (
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	timeutils "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/time"
)

type SimpleRuleChecker struct {
	rule       rules.Rule
	next       api.OutputConnector
	resultExit api.ResultConnector
}

func NewSimpleRuleChecker(r rules.Rule) *SimpleRuleChecker {
	self := new(SimpleRuleChecker)
	self.rule = r
	self.next = nil
	self.resultExit = nil

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

func (self *SimpleRuleChecker) Execute(v1, v2 visits.Visit) {
	fmt.Println("Se proceso la visita")

	durationCheck := true

	v1EntranceTime := v1.EntranceTimestamp.Time()
	v2EntranceTime := v2.EntranceTimestamp.Time()
	// v1ExitTime := v1.ExitTimestamp.Time()
	// v2ExitTime := v2.ExitTimestamp.Time()

	// Cannot get space info right now, should be part of the processed visit
	// space := v1.Space

	if self.rule.DurationCmp == "" {
		sharedTime := timeutils.AbsDateDiffInMinutes(v1EntranceTime, v2EntranceTime)
		if self.rule.DurationCmp == "<" {
			durationCheck = int(sharedTime) <= self.rule.DurationValue
		} else {
			durationCheck = int(sharedTime) >= self.rule.DurationValue
		}
	}

	if durationCheck {
		if self.resultExit != nil {
			res := api.Result{Severity: self.rule.ContagionRisk, Error: nil}
			self.resultExit.Push(res)
		}
	} else {
		if self.next != nil {
			self.next.Push(v1, v2)
		}
	}
}

func (self *SimpleRuleChecker) Push(v1, v2 visits.Visit) bool {
	self.Execute(v1, v2)
	return true
}

// -- CASTERS

func (self *SimpleRuleChecker) AsRuleChecker() api.RuleChecker {
	return self
}

func (self *SimpleRuleChecker) AsOutputConnector() api.OutputConnector {
	return self
}
