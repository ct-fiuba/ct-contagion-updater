package impl

import (
	"time"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/models/visits"
)

const DEFAULT_ESTIMATED_DURATION = 20

func GetVisitInterval(visit *visits.Visit, spaceEstimatedDuration *int) (time.Time, time.Time) {
	estimatedDuration := DEFAULT_ESTIMATED_DURATION
	if spaceEstimatedDuration != nil {
		estimatedDuration = *spaceEstimatedDuration
	}

	entranceTime := visit.EntranceTimestamp.Time()
	exitTime := entranceTime.Add(time.Minute * time.Duration(estimatedDuration))

	if visit.ExitTimestamp != nil {
		exitTime = visit.ExitTimestamp.Time()
	}

	return entranceTime, exitTime
}
