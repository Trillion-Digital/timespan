package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

func TestNewHalfMonthWindowStartingOn(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{name: "starting in first half", input: mustDate(t, "2026-03-05"), wantStart: "2026-03-05", wantEnd: "2026-03-15"},
		{name: "starting in second half mid", input: mustDate(t, "2026-03-20"), wantStart: "2026-03-20", wantEnd: "2026-03-31"},
		{name: "starting on last day of month", input: mustDate(t, "2026-03-31"), wantStart: "2026-03-31", wantEnd: "2026-03-31"},
		{name: "starting second half february", input: mustDate(t, "2026-02-20"), wantStart: "2026-02-20", wantEnd: "2026-02-28"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewHalfMonthWindowStartingOn(tt.input)
			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
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
		{name: "ending in first half", input: mustDate(t, "2026-03-10"), wantStart: "2026-03-01", wantEnd: "2026-03-10"},
		{name: "ending exactly on 15th", input: mustDate(t, "2026-03-15"), wantStart: "2026-03-01", wantEnd: "2026-03-15"},
		{name: "ending in second half", input: mustDate(t, "2026-03-20"), wantStart: "2026-03-16", wantEnd: "2026-03-20"},
		{name: "ending on last day of month", input: mustDate(t, "2026-03-31"), wantStart: "2026-03-16", wantEnd: "2026-03-31"},
		{name: "ending february last day", input: mustDate(t, "2026-02-28"), wantStart: "2026-02-16", wantEnd: "2026-02-28"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timespan.NewHalfMonthWindowEndingOn(tt.input)
			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
		})
	}
}

func TestHalfMonthWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		window    func() timespan.Window
		wantStart string
		wantEnd   string
	}{
		{
			name:      "default advances from first to second half",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-10")) },
			wantStart: "2026-03-16",
			wantEnd:   "2026-03-16",
		},
		{
			name:      "default advances from second half to next month first half",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-20")) },
			wantStart: "2026-04-01",
			wantEnd:   "2026-04-01",
		},
		{
			name:      "step month moves to same half next month",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-20")) },
			wantStart: "2026-04-16",
			wantEnd:   "2026-04-20",
		},
		{
			name:      "step year preserves end-of-month intent",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-01-31")) },
			wantStart: "2027-01-16",
			wantEnd:   "2027-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got timespan.Window
			switch tt.name {
			case "step month moves to same half next month":
				got = tt.window().Next(timespan.StepMonth)
			case "step year preserves end-of-month intent":
				got = tt.window().Next(timespan.StepYear)
			default:
				got = tt.window().Next()
			}

			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
		})
	}
}

func TestHalfMonthWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		window    func() timespan.Window
		wantStart string
		wantEnd   string
	}{
		{
			name:      "default moves from second to first half",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-20")) },
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-01",
		},
		{
			name:      "default moves from first half to previous month second half",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-10")) },
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-28",
		},
		{
			name:      "step month moves to same half previous month",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-20")) },
			wantStart: "2026-02-16",
			wantEnd:   "2026-02-20",
		},
		{
			name:      "step year preserves end-of-month intent backwards",
			window:    func() timespan.Window { return timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-31")) },
			wantStart: "2025-03-16",
			wantEnd:   "2025-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got timespan.Window
			switch tt.name {
			case "step month moves to same half previous month":
				got = tt.window().Prev(timespan.StepMonth)
			case "step year preserves end-of-month intent backwards":
				got = tt.window().Prev(timespan.StepYear)
			default:
				got = tt.window().Prev()
			}

			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
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
		{name: "complete first half", input: mustDate(t, "2026-03-10"), wantStart: "2026-03-01", wantEnd: "2026-03-15", fn: timespan.NewHalfMonthWindowEndingOn},
		{name: "complete second half", input: mustDate(t, "2026-03-20"), wantStart: "2026-03-16", wantEnd: "2026-03-31", fn: timespan.NewHalfMonthWindowEndingOn},
		{name: "starting mid half completes", input: mustDate(t, "2026-03-18"), wantStart: "2026-03-16", wantEnd: "2026-03-31", fn: timespan.NewHalfMonthWindowStartingOn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input).Complete()
			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
		})
	}
}
