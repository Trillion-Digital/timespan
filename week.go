package timespan

import "time"

type WeekWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (w *WeekWindow) Index() int {
	return weekIndex(w.end)
}

func (w *WeekWindow) Start() time.Time { return w.start }
func (w *WeekWindow) End() time.Time   { return w.end }

func (w *WeekWindow) Next(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return w.shift(12)
	}
	return w.shift(1)
}

func (w *WeekWindow) Prev(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return w.shift(-12)
	}
	return w.shift(-1)
}

func (w *WeekWindow) Complete() Window {
	ref := w.end
	if w.anchor == StartAnchor {
		ref = w.start
	}

	return &WeekWindow{
		start:           weekStart(ref),
		end:             weekEnd(ref),
		anchor:          w.anchor,
		shouldBeLastDay: true,
	}
}

func (w *WeekWindow) shift(months int) Window {
	ref := w.end
	if w.anchor == StartAnchor {
		ref = w.start
	}

	week := weekIndex(ref)

	ref = shiftMonthClamp(ref, months)

	y, m, _ := ref.Date()
	loc := ref.Location()

	switch week {
	case 0:
		ref = time.Date(y, m, ref.Day(), 0, 0, 0, 0, loc)
	case 1:
		ref = time.Date(y, m, ref.Day(), 0, 0, 0, 0, loc)
	case 2:
		ref = time.Date(y, m, ref.Day(), 0, 0, 0, 0, loc)
	case 3:
		last := time.Date(y, m+1, ref.Day(), 0, 0, 0, 0, loc).Day()
		ref = time.Date(y, m, last, 0, 0, 0, 0, loc)
	}

	if w.anchor == EndAnchor && w.shouldBeLastDay {
		ref = weekEnd(ref)
	}

	switch w.anchor {
	case StartAnchor:
		return NewWeekWindowStartingOn(ref)
	default:
		return NewWeekWindowEndingOn(ref)
	}
}

func NewWeekWindowStartingOn(t time.Time) Window {
	return &WeekWindow{
		start:           truncateToDay(t),
		end:             weekEnd(t),
		anchor:          StartAnchor,
		shouldBeLastDay: true,
	}
}

func NewWeekWindowEndingOn(t time.Time) Window {
	return &WeekWindow{
		start:           weekStart(t),
		end:             truncateToDay(t),
		anchor:          EndAnchor,
		shouldBeLastDay: isWeekEnd(t),
	}
}

func weekStart(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()

	switch {
	case d <= 7:
		return time.Date(y, m, 1, 0, 0, 0, 0, loc)
	case d <= 14:
		return time.Date(y, m, 8, 0, 0, 0, 0, loc)
	case d <= 21:
		return time.Date(y, m, 15, 0, 0, 0, 0, loc)
	default:
		return time.Date(y, m, 22, 0, 0, 0, 0, loc)
	}
}

func weekEnd(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()

	switch {
	case d <= 7:
		return time.Date(y, m, 7, 0, 0, 0, 0, loc)
	case d <= 14:
		return time.Date(y, m, 14, 0, 0, 0, 0, loc)
	case d <= 21:
		return time.Date(y, m, 21, 0, 0, 0, 0, loc)
	default:
		last := time.Date(y, m+1, 0, 0, 0, 0, 0, loc).Day()
		return time.Date(y, m, last, 0, 0, 0, 0, loc)
	}
}

func isWeekEnd(t time.Time) bool {
	return t.Equal(weekEnd(t))
}

func weekIndex(t time.Time) int {
	_, _, d := t.Date()
	switch {
	case d <= 7:
		return 0
	case d <= 14:
		return 1
	case d <= 21:
		return 2
	default:
		return 3
	}
}
