package timespan_test

import (
	"testing"
	"time"
	"timespan/timespan"
)

func TestNewQuarterWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "starting in Q1",
			input:     mustDate(t, "2026-02-10"),
			wantStart: "2026-02-10",
			wantEnd:   "2026-03-31",
		},
		{
			name:      "starting in Q2",
			input:     mustDate(t, "2026-05-20"),
			wantStart: "2026-05-20",
			wantEnd:   "2026-06-30",
		},
		{
			name:      "starting in Q3",
			input:     mustDate(t, "2026-08-01"),
			wantStart: "2026-08-01",
			wantEnd:   "2026-09-30",
		},
		{
			name:      "starting in Q4",
			input:     mustDate(t, "2026-11-15"),
			wantStart: "2026-11-15",
			wantEnd:   "2026-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewQuarterWindowStartingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}
func TestNewQuarterWindowEndingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "ending in Q1",
			input:     mustDate(t, "2026-02-10"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-02-10",
		},
		{
			name:      "ending in Q2",
			input:     mustDate(t, "2026-05-20"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-05-20",
		},
		{
			name:      "ending in Q3",
			input:     mustDate(t, "2026-08-01"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-08-01",
		},
		{
			name:      "ending in Q4",
			input:     mustDate(t, "2026-11-15"),
			wantStart: "2026-10-01",
			wantEnd:   "2026-11-15",
		},
		{
			name:      "ending on quarter boundary",
			input:     mustDate(t, "2026-06-30"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-06-30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewQuarterWindowEndingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}
func TestQuarterWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending in Q1 moves to Q2",
			input:     mustDate(t, "2026-02-10"),
			wantStart: "2026-04-01",
			wantEnd:   "2026-05-10",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "ending in Q4 moves to next year Q1",
			input:     mustDate(t, "2026-11-15"),
			wantStart: "2027-01-01",
			wantEnd:   "2027-02-15",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "starting in Q2 moves start forward",
			input:     mustDate(t, "2026-05-20"),
			wantStart: "2026-08-20",
			wantEnd:   "2026-09-30",
			fn:        timespan.NewQuarterWindowStartingOn,
		},
		{
			name:      "ending on quarter boundary moves to prev quarter end",
			input:     mustDate(t, "2026-09-30"),
			wantStart: "2026-10-01",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewQuarterWindowEndingOn,
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
func TestQuarterWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending in Q2 moves to Q1",
			input:     mustDate(t, "2026-05-20"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-02-20",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "ending in Q1 moves to previous year Q4",
			input:     mustDate(t, "2026-02-10"),
			wantStart: "2025-10-01",
			wantEnd:   "2025-11-10",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "starting in Q3 moves start back",
			input:     mustDate(t, "2026-08-01"),
			wantStart: "2026-05-01",
			wantEnd:   "2026-06-30",
			fn:        timespan.NewQuarterWindowStartingOn,
		},
		{
			name:      "ending on quarter boundary moves to prev quarter end",
			input:     mustDate(t, "2026-12-31"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-09-30",
			fn:        timespan.NewQuarterWindowEndingOn,
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
func TestQuarterWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending in Q1 completes",
			input:     mustDate(t, "2026-02-10"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "ending in Q3 completes",
			input:     mustDate(t, "2026-08-01"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-09-30",
			fn:        timespan.NewQuarterWindowEndingOn,
		},
		{
			name:      "starting in Q4 completes",
			input:     mustDate(t, "2026-11-15"),
			wantStart: "2026-10-01",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewQuarterWindowStartingOn,
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
