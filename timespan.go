package timespan

import (
	"iter"
	"time"
)

type Anchor int

const (
	StartAnchor Anchor = 0
	EndAnchor   Anchor = 1
)

type Step string

const (
	StepMonth Step = "month"
	StepYear  Step = "year"
)

func (p Step) Valid() bool {
	switch p {
	case StepMonth, StepYear:
		return true
	default:
		return false
	}
}

type Period string

const (
	Custom    Period = "custom"
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
	SetStart(t time.Time)
	Start() time.Time
	SetEnd(t time.Time)
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

func ContainsWindow(w Window, v Window) bool {
	if v.End().Before(w.Start()) {
		return false
	}

	return !v.Start().Before(w.Start()) && !v.End().After(w.End())
}

func ContainsTime(w Window, t time.Time) bool {
	return !t.Before(w.Start()) && !t.After(w.End())
}

func ContainsRange(w Window, start, end time.Time) bool {
	if end.Before(start) {
		return false
	}

	return !start.Before(w.Start()) && !end.After(w.End())
}

func truncateToDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
