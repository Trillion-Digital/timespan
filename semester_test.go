package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

func TestNewSemesterWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "starting in first semester",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-03-15",
			wantEnd:   "2026-06-30",
		},
		{
			name:      "starting in second semester",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2026-09-10",
			wantEnd:   "2026-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewSemesterWindowStartingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}
func TestNewSemesterWindowEndingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "ending in first semester",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-03-15",
		},
		{
			name:      "ending in second semester",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-09-10",
		},
		{
			name:      "ending on semester boundary",
			input:     mustDate(t, "2026-06-30"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-06-30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewSemesterWindowEndingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}
func TestSemesterWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid first semester moves to next semester",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-09-15",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "ending mid second semester moves to next year first semester",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2027-01-01",
			wantEnd:   "2027-03-10",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "starting mid first semester moves start forward",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-09-15",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewSemesterWindowStartingOn,
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
func TestSemesterWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid first semester moves to previous year second semester",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2025-07-01",
			wantEnd:   "2025-09-15",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "ending mid second semester moves to first semester",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-03-10",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "starting mid second semester moves start back",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2026-03-10",
			wantEnd:   "2026-06-30",
			fn:        timespan.NewSemesterWindowStartingOn,
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
func TestSemesterWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid first semester completes",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-06-30",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "ending mid second semester completes",
			input:     mustDate(t, "2026-09-10"),
			wantStart: "2026-07-01",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewSemesterWindowEndingOn,
		},
		{
			name:      "starting mid semester completes",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-06-30",
			fn:        timespan.NewSemesterWindowStartingOn,
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
