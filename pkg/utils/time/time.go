package time

import (
	"time"
)

func AbsDateDiffInMinutes(t1, t2 time.Time) float64 {
	if t1.After(t2) {
		return t1.Sub(t2).Minutes()
	} else {
		return t2.Sub(t1).Minutes()
	}
}

func Latest(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func Earliest(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t2
	}
	return t1
}

func IntervalsOverlap(t1_start, t1_end, t2_start, t2_end time.Time) bool {
	return !(t1_start.After(t2_end) || t2_start.After(t1_end))
}
