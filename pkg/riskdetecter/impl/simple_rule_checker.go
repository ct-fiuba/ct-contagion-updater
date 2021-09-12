package impl

import (
	"fmt"
	"time"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/compromisedCodes"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/rules"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/spaces"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
	"github.com/ct-fiuba/ct-contagion-updater/pkg/riskdetecter/api"
	timeutils "github.com/ct-fiuba/ct-contagion-updater/pkg/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimpleRuleChecker struct {
	rule       rules.Rule
	next       api.RuleChecker
	resultExit api.ResultConnector
}

func NewSimpleRuleChecker(r rules.Rule) api.RuleChecker {
	checker := new(SimpleRuleChecker)
	checker.rule = r
	checker.next = nil
	checker.resultExit = nil

	return checker
}

func (checker *SimpleRuleChecker) SetNext(rc api.RuleChecker) {
	checker.next = rc
}

func (checker *SimpleRuleChecker) SetResultExit(rc api.ResultConnector) {
	checker.resultExit = rc
}

func (checker *SimpleRuleChecker) Process(compromised, infected *visits.Visit, s *spaces.Space) error {
	durationCheck := true
	m2Check := true
	spaceCheck := true

	compromisedEntranceTime := compromised.EntranceTimestamp.Time()
	infectedEntranceTime := infected.EntranceTimestamp.Time()
	// TODO: overload with real exit time, if it exists
	compromisedExitTime := compromisedEntranceTime.Add(time.Minute * time.Duration(s.EstimatedVisitDuration))
	infectedExitTime := infectedEntranceTime.Add(time.Minute * time.Duration(s.EstimatedVisitDuration))

	fmt.Printf("[Rule #%d] Compromised time interval = [ %s ; %s ] \n", checker.rule.Index, compromisedEntranceTime.String(), compromisedExitTime.String())
	fmt.Printf("[Rule #%d] Infected time interval = [ %s ; %s ] \n", checker.rule.Index, infectedEntranceTime.String(), infectedExitTime.String())

	if checker.rule.DurationValue != 0 || !timeutils.IntervalsOverlap(compromisedEntranceTime, compromisedExitTime, infectedEntranceTime, infectedExitTime) {
		startSharedTime := timeutils.Latest(compromisedEntranceTime, infectedEntranceTime)
		endSharedTime := timeutils.Earliest(compromisedExitTime, infectedExitTime)
		sharedTime := timeutils.AbsDateDiffInMinutes(startSharedTime, endSharedTime)
		fmt.Printf("[Rule #%d] SHARED TIME INTERVAL = [ %s ; %s ] \n", checker.rule.Index, startSharedTime.String(), endSharedTime.String())
		fmt.Printf("[Rule #%d] SHARED TIME BETWEEN VISITS = %f\n", checker.rule.Index, sharedTime)
		if checker.rule.DurationCmp == "<" {
			durationCheck = int(sharedTime) <= checker.rule.DurationValue
		} else {
			durationCheck = int(sharedTime) >= checker.rule.DurationValue
		}
	}

	if checker.rule.M2Value != 0 {
		if checker.rule.M2Cmp == "<" {
			m2Check = s.M2 <= checker.rule.M2Value
		} else {
			m2Check = s.M2 >= checker.rule.M2Value
		}
	}

	// if checker.rule.SpaceValue {
	// 	spaceCheck = s.OpenPlace === rule.spaceValue;
	// }

	if durationCheck && m2Check && spaceCheck {
		fmt.Printf("[Rule #%d] Match found. Should push result to... %+v\n", checker.rule.Index, checker.resultExit)
		if checker.resultExit != nil {
			res := api.Result{CompromisedCode: compromisedCodes.CompromisedCode{
				ScanCode:          compromised.ScanCode,
				UserGeneratedCode: compromised.UserGeneratedCode,
				DateDetected:      primitive.NewDateTimeFromTime(time.Now()),
				Risk:              api.RiskStringsToSeverity[checker.rule.ContagionRisk],
			}, Error: nil}
			fmt.Printf("[Rule #%d] Match found. Pushing Result --> %+v\n", checker.rule.Index, res)
			checker.resultExit.Push(res)
		}
	} else {
		if checker.next != nil {
			checker.next.Process(compromised, infected, s)
		}
	}

	return nil
}
