package timespan

import "time"

type HalfMonthWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (h *HalfMonthWindow) Start() time.Time { return h.start }
func (h *HalfMonthWindow) End() time.Time   { return h.end }

func (h *HalfMonthWindow) Next(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return h.shift(12)
	}

	return h.shift(1)
}

func (h *HalfMonthWindow) Prev(s ...Step) Window {
	if step, ok := GetFirst(s); ok && step == StepYear {
		return h.shift(-12)
	}

	return h.shift(-1)
}

func (h *HalfMonthWindow) Complete() Window {
	ref := h.end
	if h.anchor == StartAnchor {
		ref = h.start
	}

	return &HalfMonthWindow{
		start:           halfMonthStart(ref),
		end:             halfMonthEnd(ref),
		anchor:          h.anchor,
		shouldBeLastDay: true,
	}
}

func (h *HalfMonthWindow) shift(delta int) Window {
	ref := h.end
	if h.anchor == StartAnchor {
		ref = h.start
	}

	// move one month, clamped (critical!)
	ref = shiftMonthClamp(ref, delta)

	if h.anchor == EndAnchor && h.shouldBeLastDay {
		ref = halfMonthEnd(ref)
	}

	switch h.anchor {
	case StartAnchor:
		return NewHalfMonthWindowStartingOn(ref)
	default:
		return NewHalfMonthWindowEndingOn(ref)
	}
}

func NewHalfMonthWindowStartingOn(t time.Time) Window {
	return &HalfMonthWindow{
		start:           truncateToDay(t),
		end:             halfMonthEnd(t),
		anchor:          StartAnchor,
		shouldBeLastDay: isLastDayOfMonth(t),
	}
}

func NewHalfMonthWindowEndingOn(t time.Time) Window {
	return &HalfMonthWindow{
		start:           halfMonthStart(t),
		end:             truncateToDay(t),
		anchor:          EndAnchor,
		shouldBeLastDay: isLastDayOfMonth(t),
	}
}

func halfMonthStart(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()

	if d <= 15 {
		return time.Date(y, m, 1, 0, 0, 0, 0, loc)
	}
	return time.Date(y, m, 16, 0, 0, 0, 0, loc)
}

func halfMonthEnd(t time.Time) time.Time {
	y, m, d := t.Date()
	loc := t.Location()

	if d <= 15 {
		return time.Date(y, m, 15, 0, 0, 0, 0, loc)
	}

	last := time.Date(y, m+1, 0, 0, 0, 0, 0, loc).Day()
	return time.Date(y, m, last, 0, 0, 0, 0, loc)
}

func isHalfMonthEnd(t time.Time) bool {
	return t.Equal(halfMonthEnd(t))
}
