package timespan

import "time"

type MonthlyWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (m *MonthlyWindow) Index() int {
	return int(m.end.Month())
}

func (m *MonthlyWindow) Start() time.Time { return m.start }
func (m *MonthlyWindow) End() time.Time   { return m.end }

// ----- core navigation -----

func (m *MonthlyWindow) Next(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return m.shift(12)
	}

	return m.shift(1)
}

func (m *MonthlyWindow) Prev(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return m.shift(-12)
	}

	return m.shift(-1)
}

func (m *MonthlyWindow) shift(delta int) Window {
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
		return NewMonthlyWindowStartingOn(ref)
	default:
		return NewMonthlyWindowEndingOn(ref)
	}
}

// ----- projection -----

func (m *MonthlyWindow) Complete() Window {
	y, mo, _ := m.end.Date()
	loc := m.end.Location()

	start := time.Date(y, mo, 1, 0, 0, 0, 0, loc)
	end := time.Date(y, mo+1, 0, 0, 0, 0, 0, loc)

	return &MonthlyWindow{
		start:           start,
		end:             end,
		anchor:          m.anchor,
		shouldBeLastDay: true,
	}
}

// ----- constructors -----

func NewMonthlyWindowStartingOn(t time.Time) Window {
	y, m, _ := t.Date()
	loc := t.Location()

	end := time.Date(y, m+1, 0, 0, 0, 0, 0, loc)

	return &MonthlyWindow{
		start:           truncateToDay(t),
		end:             end,
		anchor:          StartAnchor,
		shouldBeLastDay: isLastDayOfMonth(t),
	}
}

func NewMonthlyWindowEndingOn(t time.Time) Window {
	y, m, _ := t.Date()
	loc := t.Location()

	start := time.Date(y, m, 1, 0, 0, 0, 0, loc)

	return &MonthlyWindow{
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
