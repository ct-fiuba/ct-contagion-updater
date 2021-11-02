package impl

import (
	"time"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

func GetVisitInterval(visit *visits.Visit, spaceEstimatedDuration int) (time.Time, time.Time) {
	entranceTime := visit.EntranceTimestamp.Time()
	exitTime := entranceTime.Add(time.Minute * time.Duration(spaceEstimatedDuration))

	if visit.ExitTimestamp != nil {
		exitTime = visit.ExitTimestamp.Time()
	}

	return entranceTime, exitTime
}
