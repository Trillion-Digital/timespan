package timespan

import "time"

type QuarterWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (q *QuarterWindow) Index() int {
	_, m, _ := q.end.Date()

	return int(m / 3)
}

func (q *QuarterWindow) Start() time.Time { return q.start }
func (q *QuarterWindow) End() time.Time   { return q.end }

func (q *QuarterWindow) Next(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return q.shift(12)
	}

	return q.shift(3)
}

func (q *QuarterWindow) Prev(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return q.shift(-12)
	}

	return q.shift(-3)
}

func (q *QuarterWindow) Complete() Window {
	ref := q.end
	if q.anchor == StartAnchor {
		ref = q.start
	}

	return &QuarterWindow{
		start:           quarterStart(ref),
		end:             quarterEnd(ref),
		anchor:          q.anchor,
		shouldBeLastDay: true,
	}
}

func (q *QuarterWindow) shift(months int) Window {
	ref := q.end
	if q.anchor == StartAnchor {
		ref = q.start
	}

	ref = shiftMonthClamp(ref, months)

	if q.anchor == EndAnchor && q.shouldBeLastDay {
		ref = quarterEnd(ref)
	}

	if q.anchor == StartAnchor {
		return NewQuarterWindowStartingOn(ref)
	}
	return NewQuarterWindowEndingOn(ref)
}

func NewQuarterWindowStartingOn(t time.Time) Window {
	return &QuarterWindow{
		start:           truncateToDay(t),
		end:             quarterEnd(t),
		anchor:          StartAnchor,
		shouldBeLastDay: true,
	}
}

func NewQuarterWindowEndingOn(t time.Time) Window {
	return &QuarterWindow{
		start:           quarterStart(t),
		end:             truncateToDay(t),
		anchor:          EndAnchor,
		shouldBeLastDay: isQuarterEnd(t),
	}
}

func quarterStart(t time.Time) time.Time {
	y, m, _ := t.Date()
	loc := t.Location()

	switch {
	case m <= 3:
		return time.Date(y, 1, 1, 0, 0, 0, 0, loc)
	case m <= 6:
		return time.Date(y, 4, 1, 0, 0, 0, 0, loc)
	case m <= 9:
		return time.Date(y, 7, 1, 0, 0, 0, 0, loc)
	default:
		return time.Date(y, 10, 1, 0, 0, 0, 0, loc)
	}
}

func quarterEnd(t time.Time) time.Time {
	y, m, _ := t.Date()
	loc := t.Location()

	switch {
	case m <= 3:
		return time.Date(y, 3, 31, 0, 0, 0, 0, loc)
	case m <= 6:
		return time.Date(y, 6, 30, 0, 0, 0, 0, loc)
	case m <= 9:
		return time.Date(y, 9, 30, 0, 0, 0, 0, loc)
	default:
		return time.Date(y, 12, 31, 0, 0, 0, 0, loc)
	}
}

func isQuarterEnd(t time.Time) bool {
	return t.Equal(quarterEnd(t))
}
