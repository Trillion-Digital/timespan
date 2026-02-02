package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

func TestNewYearlyWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "starting mid year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-03-15",
			wantEnd:   "2026-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewYearlyWindowStartingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestNewYearlyWindowEndingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "ending mid year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-03-15",
		},
		{
			name:      "ending last day of year",
			input:     mustDate(t, "2026-12-31"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewYearlyWindowEndingOn(tt.input)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestYearlyWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid year moves to next year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2027-01-01",
			wantEnd:   "2027-03-15",
			fn:        timespan.NewYearlyWindowEndingOn,
		},
		{
			name:      "ending last day moves to full next year",
			input:     mustDate(t, "2026-12-31"),
			wantStart: "2027-01-01",
			wantEnd:   "2027-12-31",
			fn:        timespan.NewYearlyWindowEndingOn,
		},
		{
			name:      "starting mid year moves start forward",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2027-03-15",
			wantEnd:   "2027-12-31",
			fn:        timespan.NewYearlyWindowStartingOn,
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

func TestYearlyWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid year moves back one year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2025-01-01",
			wantEnd:   "2025-03-15",
			fn:        timespan.NewYearlyWindowEndingOn,
		},
		{
			name:      "ending last day moves to full previous year",
			input:     mustDate(t, "2026-12-31"),
			wantStart: "2025-01-01",
			wantEnd:   "2025-12-31",
			fn:        timespan.NewYearlyWindowEndingOn,
		},
		{
			name:      "starting mid year moves start back",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2025-03-15",
			wantEnd:   "2025-12-31",
			fn:        timespan.NewYearlyWindowStartingOn,
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

func TestYearlyWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "ending mid year completes to full year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewYearlyWindowEndingOn,
		},
		{
			name:      "starting mid year completes to full year",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-12-31",
			fn:        timespan.NewYearlyWindowStartingOn,
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
