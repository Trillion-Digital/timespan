package timespan

import "time"

type SemesterWindow struct {
	start           time.Time
	end             time.Time
	anchor          Anchor
	shouldBeLastDay bool
}

func (s *SemesterWindow) Index() int {
	_, m, _ := s.end.Date()

	if m <= 6 {
		return 0
	}
	return 1
}

func (s *SemesterWindow) Start() time.Time { return s.start }
func (s *SemesterWindow) End() time.Time   { return s.end }

func NewSemesterWindowStartingOn(t time.Time) Window {
	return &SemesterWindow{
		start:           truncateToDay(t),
		end:             semesterEnd(t),
		anchor:          StartAnchor,
		shouldBeLastDay: isSemesterEnd(t),
	}
}

func NewSemesterWindowEndingOn(t time.Time) Window {
	return &SemesterWindow{
		start:           semesterStart(t),
		end:             truncateToDay(t),
		anchor:          EndAnchor,
		shouldBeLastDay: isSemesterEnd(t),
	}
}

func (s *SemesterWindow) Next(st ...Step) Window {
	if step, ok := GetFirst(st); ok && step == StepYear {
		return s.shift(12)
	}

	return s.shift(6)
}

func (s *SemesterWindow) Prev(st ...Step) Window {
	if step, ok := GetFirst(st); ok && step == StepYear {
		return s.shift(-12)
	}

	return s.shift(-6)
}

func (s *SemesterWindow) shift(months int) Window {
	ref := s.end
	if s.anchor == StartAnchor {
		ref = s.start
	}

	ref = ref.AddDate(0, months, 0)

	if s.shouldBeLastDay {
		ref = semesterEnd(ref)
	}

	switch s.anchor {
	case StartAnchor:
		return NewSemesterWindowStartingOn(ref)
	default:
		return NewSemesterWindowEndingOn(ref)
	}
}

func (s *SemesterWindow) Complete() Window {
	ref := s.end
	if s.anchor == StartAnchor {
		ref = s.start
	}

	return &SemesterWindow{
		start:           semesterStart(ref),
		end:             semesterEnd(ref),
		anchor:          s.anchor,
		shouldBeLastDay: true,
	}
}

func semesterStart(t time.Time) time.Time {
	y, m, _ := t.Date()
	loc := t.Location()

	if m <= 6 {
		return time.Date(y, 1, 1, 0, 0, 0, 0, loc)
	}
	return time.Date(y, 7, 1, 0, 0, 0, 0, loc)
}

func semesterEnd(t time.Time) time.Time {
	y, m, _ := t.Date()
	loc := t.Location()

	if m <= 6 {
		return time.Date(y, 6, 30, 0, 0, 0, 0, loc)
	}
	return time.Date(y, 12, 31, 0, 0, 0, 0, loc)
}

func isSemesterEnd(t time.Time) bool {
	return t.Equal(semesterEnd(t))
}
