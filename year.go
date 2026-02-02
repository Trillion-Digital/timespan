package timespan

import "time"

type YearlyWindow struct {
	start  time.Time
	end    time.Time
	anchor Anchor
}

func (y *YearlyWindow) Start() time.Time { return y.start }
func (y *YearlyWindow) End() time.Time   { return y.end }

func NewYearlyWindowStartingOn(t time.Time) Window {
	y, _, _ := t.Date()
	loc := t.Location()

	end := time.Date(y, 12, 31, 0, 0, 0, 0, loc)

	return &YearlyWindow{
		start:  truncateToDay(t),
		end:    end,
		anchor: StartAnchor,
	}
}

func NewYearlyWindowEndingOn(t time.Time) Window {
	y, _, _ := t.Date()
	loc := t.Location()

	start := time.Date(y, 1, 1, 0, 0, 0, 0, loc)

	return &YearlyWindow{
		start:  start,
		end:    truncateToDay(t),
		anchor: EndAnchor,
	}
}

func (y *YearlyWindow) Next(s ...Step) Window {
	return y.shift(1)
}

func (y *YearlyWindow) Prev(s ...Step) Window {
	return y.shift(-1)
}

func (y *YearlyWindow) Complete() Window {
	year, _, _ := y.end.Date()
	loc := y.end.Location()

	start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(year, 12, 31, 0, 0, 0, 0, loc)

	return &YearlyWindow{
		start:  start,
		end:    end,
		anchor: y.anchor,
	}
}

func (y *YearlyWindow) shift(delta int) Window {
	ref := y.end
	if y.anchor == StartAnchor {
		ref = y.start
	}

	ref = ref.AddDate(delta, 0, 0)

	switch y.anchor {
	case StartAnchor:
		return NewYearlyWindowStartingOn(ref)
	default:
		return NewYearlyWindowEndingOn(ref)
	}
}
