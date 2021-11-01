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
	executedChecks := make(map[string]bool)

	if checker.rule.DurationValue != nil {
		compromisedEntranceTime, compromisedExitTime := GetVisitInterval(compromised, s.EstimatedVisitDuration)
		infectedEntranceTime, infectedExitTime := GetVisitInterval(infected, s.EstimatedVisitDuration)
		startSharedTime := timeutils.Latest(compromisedEntranceTime, infectedEntranceTime)
		endSharedTime := timeutils.Earliest(compromisedExitTime, infectedExitTime)
		sharedTime := timeutils.AbsDateDiffInMinutes(startSharedTime, endSharedTime)

		durationCheck := false
		if *checker.rule.DurationCmp == "<" {
			durationCheck = int(sharedTime) <= *checker.rule.DurationValue
		} else {
			durationCheck = int(sharedTime) >= *checker.rule.DurationValue
		}
		executedChecks["durationCheck"] = durationCheck
	}

	if checker.rule.M2Value != nil {
		m2Check := false
		if *checker.rule.M2Cmp == "<" {
			m2Check = s.M2 <= *checker.rule.M2Value
		} else {
			m2Check = s.M2 >= *checker.rule.M2Value
		}
		executedChecks["m2Check"] = m2Check
	}

	if checker.rule.OpenSpace != nil {
		executedChecks["spaceCheck"] = s.OpenSpace == *checker.rule.OpenSpace
	}

	if checker.rule.N95Mandatory != nil {
		executedChecks["n95MandatoryCheck"] = s.N95Mandatory == *checker.rule.N95Mandatory
	}

	if checker.rule.Vaccinated != nil && compromised.Vaccinated != nil {
		executedChecks["vaccinatedCheck"] = *compromised.Vaccinated == *checker.rule.Vaccinated
	}

	if checker.rule.VaccineReceived != nil && compromised.VaccineReceived != nil {
		executedChecks["vaccineReceivedCheck"] = *compromised.VaccineReceived == *checker.rule.VaccineReceived
	}

	if checker.rule.VaccinatedDaysAgoMin != nil && compromised.VaccinatedDate != nil {
		vaccineDate := (*compromised.VaccinatedDate).Time()
		executedChecks["vaccinatedDaysCheck"] = int(time.Since(vaccineDate).Hours())/24 >= *checker.rule.VaccinatedDaysAgoMin
	}

	if checker.rule.IllnessRecovered != nil && compromised.IllnessRecovered != nil {
		executedChecks["illnessRecoveredCheck"] = *compromised.IllnessRecovered == *checker.rule.IllnessRecovered
	}

	if checker.rule.IllnessRecoveredDaysAgoMax != nil && compromised.IllnessRecoveredDate != nil {
		recoveredDate := (*compromised.IllnessRecoveredDate).Time()
		executedChecks["illnessRecoveredDaysCheck"] = int(time.Since(recoveredDate).Hours())/24 <= *checker.rule.IllnessRecoveredDaysAgoMax
	}

	// -- Decision
	failedChecks := 0
	for check, result := range executedChecks {
		fmt.Printf("[Rule #%d]   - %s --> %t\n", checker.rule.Index, check, result)
		if result {
			failedChecks = failedChecks + 1
		}
	}

	if failedChecks == len(executedChecks) && failedChecks != 0 {
		fmt.Printf("[Rule #%d] Match found between visits %s and %s\n", checker.rule.Index, compromised.UserGeneratedCode, infected.UserGeneratedCode)
		if checker.resultExit != nil {
			res := api.Result{CompromisedCode: compromisedCodes.CompromisedCode{
				SpaceId:           compromised.SpaceId,
				UserGeneratedCode: compromised.UserGeneratedCode,
				DateDetected:      primitive.NewDateTimeFromTime(time.Now()),
				Risk:              checker.rule.ContagionRisk,
			}, Error: nil}
			checker.resultExit.Push(res)
		}
	} else {
		if checker.next != nil {
			checker.next.Process(compromised, infected, s)
		}
	}

	return nil
}
