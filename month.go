package timespan

import "time"

type MonthWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (m *MonthWindow) Index() int {
	return int(m.end.Month())
}

func (m *MonthWindow) Start() time.Time { return m.start }
func (m *MonthWindow) End() time.Time   { return m.end }

// ----- core navigation -----

func (m *MonthWindow) Next(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return m.shift(12)
	}

	return m.shift(1)
}

func (m *MonthWindow) Prev(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return m.shift(-12)
	}

	return m.shift(-1)
}

func (m *MonthWindow) shift(delta int) Window {
	ref := m.end
	if m.anchor == StartAnchor {
		ref = m.start
	}

	ref = shiftMonthClamp(ref, delta)

	if m.shouldBeLastDay {
		ref = snapToLastDayOfMonth(ref)
	}

	switch m.anchor {
	case StartAnchor:
		return NewMonthWindowStartingOn(ref)
	default:
		return NewMonthWindowEndingOn(ref)
	}
}

// ----- projection -----

func (m *MonthWindow) Complete() Window {
	y, mo, _ := m.end.Date()
	loc := m.end.Location()

	start := time.Date(y, mo, 1, 0, 0, 0, 0, loc)
	end := time.Date(y, mo+1, 0, 0, 0, 0, 0, loc)

	return &MonthWindow{
		start:           start,
		end:             end,
		anchor:          m.anchor,
		shouldBeLastDay: true,
	}
}

// ----- constructors -----

func NewMonthWindowStartingOn(t time.Time) Window {
	y, m, _ := t.Date()
	loc := t.Location()

	end := time.Date(y, m+1, 0, 0, 0, 0, 0, loc)

	return &MonthWindow{
		start:           truncateToDay(t),
		end:             end,
		anchor:          StartAnchor,
		shouldBeLastDay: isLastDayOfMonth(t),
	}
}

func NewMonthWindowEndingOn(t time.Time) Window {
	y, m, _ := t.Date()
	loc := t.Location()

	start := time.Date(y, m, 1, 0, 0, 0, 0, loc)

	return &MonthWindow{
		start:           start,
		end:             truncateToDay(t),
		anchor:          EndAnchor,
		shouldBeLastDay: isLastDayOfMonth(t),
	}
}

// ----- date helpers -----

func shiftMonthClamp(t time.Time, delta int) time.Time {
	y, m, d := t.Date()
	loc := t.Location()

	first := time.Date(y, m+time.Month(delta), 1, 0, 0, 0, 0, loc)
	last := time.Date(first.Year(), first.Month()+1, 0, 0, 0, 0, 0, loc).Day()

	if d > last {
		d = last
	}

	return time.Date(first.Year(), first.Month(), d, 0, 0, 0, 0, loc)
}

func snapToLastDayOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	loc := t.Location()

	last := time.Date(y, m+1, 0, 0, 0, 0, 0, loc).Day()
	return time.Date(y, m, last, 0, 0, 0, 0, loc)
}

func isLastDayOfMonth(t time.Time) bool {
	y, m, d := t.Date()
	last := time.Date(y, m+1, 0, 0, 0, 0, 0, t.Location()).Day()
	return d == last
}
