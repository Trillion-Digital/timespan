package timespan

import (
	"iter"
	"time"
)

type Step int

const (
	StepPeriod Step = iota
	StepYear
)

type Anchor int

const (
	StartAnchor Anchor = 0
	EndAnchor   Anchor = 1
)

type Period string

const (
	Week      Period = "week"
	HalfMonth Period = "halfmonth"
	Month     Period = "month"
	Quarter   Period = "quarter"
	Semester  Period = "semester"
	Year      Period = "year"
)

func (p Period) Valid() bool {
	switch p {
	case Week, HalfMonth, Month, Quarter, Semester, Year:
		return true
	default:
		return false
	}
}

type Window interface {
	Start() time.Time
	End() time.Time
	Complete() Window
	Next(s ...Step) Window
	Prev(s ...Step) Window
	Index() int
}

func WindowEndingOn(period Period, t time.Time) Window {
	return nil
}

func WindowStartingOn(period Period, t time.Time) Window {
	return nil
}

func Days(w Window) iter.Seq[time.Time] {
	start := truncateToDay(w.Start())
	end := truncateToDay(w.End())

	return func(yield func(time.Time) bool) {
		for !start.After(end) {
			if !yield(start) {
				return
			}

			start = start.AddDate(0, 0, 1)
		}
	}
}

func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
