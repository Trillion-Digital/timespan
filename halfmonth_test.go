package timespan_test

import (
	"testing"
	"time"
	"timespan/timespan"
)

func TestNewHalfMonthWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "starting in first half",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-03-05",
			wantEnd:   "2026-03-15",
		},
		{
			name:      "starting in second half mid",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-03-20",
			wantEnd:   "2026-03-31",
		},
		{
			name:      "starting on last day of month",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-03-31",
			wantEnd:   "2026-03-31",
		},
		{
			name:      "starting second half february",
			input:     mustDate(t, "2026-02-20"),
			wantStart: "2026-02-20",
			wantEnd:   "2026-02-28",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewHalfMonthWindowStartingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestNewHalfMonthWindowEndingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "ending in first half",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-10",
		},
		{
			name:      "ending exactly on 15th",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-15",
		},
		{
			name:      "ending in second half",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-03-16",
			wantEnd:   "2026-03-20",
		},
		{
			name:      "ending on last day of month",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-03-16",
			wantEnd:   "2026-03-31",
		},
		{
			name:      "ending february last day",
			input:     mustDate(t, "2026-02-28"),
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-28",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewHalfMonthWindowEndingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestHalfMonthWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "first half moves to first half next month",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-04-10",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "second half moves to second half next month",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-04-16",
			wantEnd:   "2026-04-20",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "ending on last day preserves last-day intent",
			input:     mustDate(t, "2026-01-31"),
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "starting in first half moves start forward",
			input:     mustDate(t, "2026-03-05"),
			wantStart: "2026-04-05",
			wantEnd:   "2026-04-15",
			fn:        timespan.NewHalfMonthWindowStartingOn,
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

func TestHalfMonthWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "second half moves to second half prev month",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-20",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "first half moves to previous month first half",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-02-01",
			wantEnd:   "2026-02-10",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "ending on last day preserves last-day intent backwards",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "starting in second half moves start back",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-02-20",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewHalfMonthWindowStartingOn,
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

func TestHalfMonthWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "complete first half",
			input:     mustDate(t, "2026-03-10"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-15",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "complete second half",
			input:     mustDate(t, "2026-03-20"),
			wantStart: "2026-03-16",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewHalfMonthWindowEndingOn,
		},
		{
			name:      "starting mid half completes",
			input:     mustDate(t, "2026-03-18"),
			wantStart: "2026-03-16",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewHalfMonthWindowStartingOn,
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
