package policy

import "time"

const (
	MaxScaleReplicas   = 5
	CountRequestPeriod = time.Minute * 1
)

// most easy scale policy
func CalculateScale(calls int, actualReplicas int) (int, bool) {
	// average calls per minute
	avg := (int)(time.Minute/CountRequestPeriod) * calls

	// so stupid
	var target int
	if avg == 0 {
		target = 0
	} else if avg < 10 {
		target = 1
	} else if avg < 20 {
		target = 2
	} else if avg < 30 {
		target = 3
	} else if avg < 40 {
		target = 4
	} else {
		target = 5
	}

	if actualReplicas < target {
		return actualReplicas + 1, true
	} else if actualReplicas > target {
		return actualReplicas - 1, true
	} else {
		return actualReplicas, false
	}
}
