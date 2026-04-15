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
		useStep   bool
		step      timespan.Step
	}{
		{
			name:      "default step moves one week for ending anchor",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-03-08",
			wantEnd:   "2026-03-12",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "month step moves week 1 to next month",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-04-05",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "month step moves week 3 to next month",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-04-15",
			wantEnd:   "2026-04-18",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "month step preserves last-day intent",
			input:     mustDate(t, "2026-01-31"),
			wantStart: "2026-02-22",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "default step moves one week for starting anchor",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-03-17",
			wantEnd:   "2026-03-21",
			fn:        timespan.NewWeekWindowStartingOn,
		},
		{
			name:      "month step moves starting anchor forward",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-04-10",
			wantEnd:   "2026-04-14",
			fn:        timespan.NewWeekWindowStartingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "year step moves to same week next year",
			input:     mustDate(t, "2026-03-14"),
			wantStart: "2027-03-08",
			wantEnd:   "2027-03-14",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fn(tt.input)

			var got timespan.Window
			if tt.useStep {
				got = w.Next(tt.step)
			} else {
				got = w.Next()
			}

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
		useStep   bool
		step      timespan.Step
	}{
		{
			name:      "default step moves one week back for ending anchor",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-03",
			fn:        timespan.NewWeekWindowEndingOn,
		},
		{
			name:      "month step moves week 2 to previous month",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-02-08",
			wantEnd:   "2026-02-10",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "month step preserves last-day backwards",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-02-22",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "default step moves one week back for starting anchor",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-03-11",
			wantEnd:   "2026-03-14",
			fn:        timespan.NewWeekWindowStartingOn,
		},
		{
			name:      "month step moves starting anchor back",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-02-18",
			wantEnd:   "2026-02-21",
			fn:        timespan.NewWeekWindowStartingOn,
			useStep:   true,
			step:      timespan.StepMonth,
		},
		{
			name:      "year step moves to same week previous year",
			input:     mustDate(t, "2026-03-14"),
			wantStart: "2025-03-08",
			wantEnd:   "2025-03-14",
			fn:        timespan.NewWeekWindowEndingOn,
			useStep:   true,
			step:      timespan.StepYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.fn(tt.input)

			var got timespan.Window
			if tt.useStep {
				got = w.Prev(tt.step)
			} else {
				got = w.Prev()
			}

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
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
