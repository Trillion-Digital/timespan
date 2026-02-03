package timespan

import "time"

type YearWindow struct {
	start  time.Time
	end    time.Time
	anchor Anchor
}

func (y *YearWindow) Index() int {
	return 1
}

func (y *YearWindow) Start() time.Time { return y.start }
func (y *YearWindow) SetStart(t time.Time) {
	y.start = t
}
func (y *YearWindow) End() time.Time { return y.end }
func (y *YearWindow) SetEnd(t time.Time) {
	y.end = t
}

func (y *YearWindow) Next(s ...Step) Window {
	return y.shift(1)
}

func (y *YearWindow) Prev(s ...Step) Window {
	return y.shift(-1)
}

func (y *YearWindow) Complete() Window {
	year, _, _ := y.end.Date()
	loc := y.end.Location()

	start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(year, 12, 31, 0, 0, 0, 0, loc)

	return &YearWindow{
		start:  start,
		end:    end,
		anchor: y.anchor,
	}
}

func (y *YearWindow) shift(delta int) Window {
	ref := y.end
	if y.anchor == StartAnchor {
		ref = y.start
	}

	ref = ref.AddDate(delta, 0, 0)

	switch y.anchor {
	case StartAnchor:
		return NewYearWindowStartingOn(ref)
	default:
		return NewYearWindowEndingOn(ref)
	}
}

func NewYearWindowStartingOn(t time.Time) Window {
	y, _, _ := t.Date()
	loc := t.Location()

	end := time.Date(y, 12, 31, 0, 0, 0, 0, loc)

	return &YearWindow{
		start:  truncateToDay(t),
		end:    end,
		anchor: StartAnchor,
	}
}

func NewYearWindowEndingOn(t time.Time) Window {
	y, _, _ := t.Date()
	loc := t.Location()

	start := time.Date(y, 1, 1, 0, 0, 0, 0, loc)

	return &YearWindow{
		start:  start,
		end:    truncateToDay(t),
		anchor: EndAnchor,
	}
}
