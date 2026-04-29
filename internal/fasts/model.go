package fasts

import "time"

type fastRecord struct {
	ID          int
	UserID      int
	StartTime   time.Time
	EndTime     *time.Time
	GoalSeconds int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Fast struct {
	fastRecord
	IsActive         bool
	ElapsedSeconds   int
	ExpectedEnd      time.Time
	RemainingSeconds int
	IsGoalReached    bool
}

func compute(f fastRecord) Fast {
	now := time.Now()
	endRef := now
	if f.EndTime != nil {
		endRef = *f.EndTime
	}

	expected := f.StartTime.Add(time.Duration(f.GoalSeconds) * time.Second)

	elapsed := max(endRef.Sub(f.StartTime), 0)

	remaining := max(expected.Sub(now), 0)

	return Fast{
		fastRecord:       f,
		IsActive:         f.EndTime == nil,
		ElapsedSeconds:   int(elapsed.Seconds()),
		ExpectedEnd:      expected,
		RemainingSeconds: int(remaining.Seconds()),
		IsGoalReached:    !now.Before(expected),
	}
}
