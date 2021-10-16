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
	n95MandatoryCheck := true
	vaccinatedCheck := true
	vaccineReceivedCheck := true
	vaccinatedDaysCheck := true
	illnessRecoveredCheck := true
	illnessRecoveredDaysCheck := true

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

	if checker.rule.OpenSpace {
		spaceCheck = s.OpenSpace == checker.rule.OpenSpace
	}

	if checker.rule.N95Mandatory {
		n95MandatoryCheck = s.N95Mandatory == checker.rule.N95Mandatory
	}

	if checker.rule.Vaccinated != 0 {
		vaccinatedCheck = compromised.Vaccinated == checker.rule.Vaccinated
	}

	if checker.rule.VaccineReceived != "" {
		vaccineReceivedCheck = compromised.VaccineReceived == checker.rule.VaccineReceived
	}

	if checker.rule.VaccinatedDaysAgoMin != 0 {
		vaccineDate := compromised.VaccinatedDate.Time()
		vaccinatedDaysCheck = int(time.Since(vaccineDate).Hours())/24 >= checker.rule.VaccinatedDaysAgoMin
	}

	if checker.rule.IllnessRecovered {
		illnessRecoveredCheck = compromised.IllnessRecovered == checker.rule.IllnessRecovered
	}

	if checker.rule.IllnessRecoveredDaysAgoMax != 0 {
		recoveredDate := compromised.IllnessRecoveredDate.Time()
		illnessRecoveredDaysCheck = int(time.Since(recoveredDate).Hours())/24 <= checker.rule.IllnessRecoveredDaysAgoMax
	}

	// -- Decision
	if durationCheck && m2Check && spaceCheck && n95MandatoryCheck && vaccinatedCheck &&
		vaccineReceivedCheck && vaccinatedDaysCheck && illnessRecoveredCheck && illnessRecoveredDaysCheck {
		fmt.Printf("[Rule #%d] Match found between visits %d and %d\n", checker.rule.Index, compromised.UserGeneratedCode, infected.UserGeneratedCode)
		if checker.resultExit != nil {
			res := api.Result{CompromisedCode: compromisedCodes.CompromisedCode{
				SpaceId:          compromised.SpaceId,
				UserGeneratedCode: compromised.UserGeneratedCode,
				DateDetected:      primitive.NewDateTimeFromTime(time.Now()),
				Risk:              api.RiskStringsToSeverity[checker.rule.ContagionRisk],
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
