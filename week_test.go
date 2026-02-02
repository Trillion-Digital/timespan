package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

func TestNewWeekWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "starting in week 1",
			input:     mustDate(t, "2026-03-03"),
			wantStart: "2026-03-03",
			wantEnd:   "2026-03-07",
		},
		{
			name:      "starting in week 2",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-03-10",
			wantEnd:   "2026-03-14",
		},
		{
			name:      "starting in week 3",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-03-18",
			wantEnd:   "2026-03-21",
		},
		{
			name:      "starting in week 4",
			input:     mustDate(t, "2026-03-25"),
			wantStart: "2026-03-25",
			wantEnd:   "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewWeekWindowStartingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestNewWeekWindowEndingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "ending in week 1",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-05",
		},
		{
			name:      "ending in week 2",
			input:     mustDate(t, "2026-03-14"),
			wantStart: "2026-03-08",
			wantEnd:   "2026-03-14",
		},
		{
			name:      "ending in week 3",
			input:     mustDate(t, "2026-03-21"),
			wantStart: "2026-03-15",
			wantEnd:   "2026-03-21",
		},
		{
			name:      "ending in week 4",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-03-22",
			wantEnd:   "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewWeekWindowEndingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestWeekWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "week 1 moves to week 1 next month",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-04-05",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "week 3 moves to week 3 next month",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-04-15",
			wantEnd:   "2026-04-18",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "week 4 preserves last-day intent",
			input:     mustDate(t, "2026-01-31"),
			wantStart: "2026-02-22",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "starting anchor moves start forward",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-04-10",
			wantEnd:   "2026-04-14",
			fn:        timespan.NewWeekWindowStartingOn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fn(tt.input)
			got := w.Next()

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestWeekWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "week 2 moves to week 2 prev month",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-02-08",
			wantEnd:   "2026-02-10",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "week 4 preserves last-day backwards",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-02-22",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "starting anchor moves start back",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-02-18",
			wantEnd:   "2026-02-21",
			fn:        timespan.NewWeekWindowStartingOn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fn(tt.input)
			got := w.Prev()

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestWeekWindow_Next_Year(t *testing.T) {
	w := timespan.NewWeekWindowEndingOn(mustDate(t, "2026-03-14"))

	got := w.Next(timespan.StepYear)

	assertWindow(
		t,
		got,
		mustDate(t, "2027-03-08"),
		mustDate(t, "2027-03-14"),
	)
}

func TestWeekWindow_Prev_Year(t *testing.T) {
	w := timespan.NewWeekWindowEndingOn(mustDate(t, "2026-03-14"))

	got := w.Prev(timespan.StepYear)

	assertWindow(
		t,
		got,
		mustDate(t, "2025-03-08"),
		mustDate(t, "2025-03-14"),
	)
}

func TestWeekWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "complete week 1",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-07",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "complete week 4",
			input:     mustDate(t, "2026-03-28"),
			wantStart: "2026-03-22",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewWeekWindowStartingOn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fn(tt.input)
			got := w.Complete()

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}
