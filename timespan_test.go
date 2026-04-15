package timespan_test

import (
	"testing"

	"github.com/Trillion-Digital/timespan"
)

func TestStepValid(t *testing.T) {
	tests := []struct {
		name string
		step timespan.Step
		want bool
	}{
		{name: "month", step: timespan.StepMonth, want: true},
		{name: "year", step: timespan.StepYear, want: true},
		{name: "empty", step: "", want: false},
		{name: "unsupported", step: "week", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.step.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodValid(t *testing.T) {
	tests := []struct {
		name   string
		period timespan.Period
		want   bool
	}{
		{name: "week", period: timespan.Week, want: true},
		{name: "half month", period: timespan.HalfMonth, want: true},
		{name: "month", period: timespan.Month, want: true},
		{name: "quarter", period: timespan.Quarter, want: true},
		{name: "semester", period: timespan.Semester, want: true},
		{name: "year", period: timespan.Year, want: true},
		{name: "custom is currently invalid", period: timespan.Custom, want: false},
		{name: "unknown", period: "nope", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.period.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDays(t *testing.T) {
	w := timespan.NewCustomWindow(mustDate(t, "2026-01-10"), mustDate(t, "2026-01-12"))
	days := collectDays(timespan.Days(w))

	if len(days) != 3 {
		t.Fatalf("len(days) = %d, want 3", len(days))
	}
	if !days[0].Equal(mustDate(t, "2026-01-10")) || !days[1].Equal(mustDate(t, "2026-01-11")) || !days[2].Equal(mustDate(t, "2026-01-12")) {
		t.Fatalf("unexpected days: %v", days)
	}
}

func TestContainsTime(t *testing.T) {
	w := timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-20"))
	tests := []struct {
		name string
		time string
		want bool
	}{
		{name: "at start", time: "2026-01-01", want: true},
		{name: "inside", time: "2026-01-10", want: true},
		{name: "at end", time: "2026-01-20", want: true},
		{name: "before", time: "2025-12-31", want: false},
		{name: "after", time: "2026-01-21", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timespan.ContainsTime(w, mustDate(t, tt.time)); got != tt.want {
				t.Fatalf("ContainsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsRange(t *testing.T) {
	w := timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-20"))
	tests := []struct {
		name  string
		start string
		end   string
		want  bool
	}{
		{name: "exact bounds", start: "2026-01-01", end: "2026-01-20", want: true},
		{name: "inside bounds", start: "2026-01-05", end: "2026-01-10", want: true},
		{name: "starts before", start: "2025-12-31", end: "2026-01-10", want: false},
		{name: "ends after", start: "2026-01-10", end: "2026-01-21", want: false},
		{name: "inverted range", start: "2026-01-10", end: "2026-01-09", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timespan.ContainsRange(w, mustDate(t, tt.start), mustDate(t, tt.end)); got != tt.want {
				t.Fatalf("ContainsRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsWindow(t *testing.T) {
	w := timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-20"))
	tests := []struct {
		name  string
		inner timespan.Window
		want  bool
	}{
		{name: "exactly same", inner: timespan.NewCustomWindow(mustDate(t, "2026-01-01"), mustDate(t, "2026-01-20")), want: true},
		{name: "fully inside", inner: timespan.NewCustomWindow(mustDate(t, "2026-01-05"), mustDate(t, "2026-01-10")), want: true},
		{name: "starts before", inner: timespan.NewCustomWindow(mustDate(t, "2025-12-31"), mustDate(t, "2026-01-10")), want: false},
		{name: "ends after", inner: timespan.NewCustomWindow(mustDate(t, "2026-01-10"), mustDate(t, "2026-01-21")), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timespan.ContainsWindow(w, tt.inner); got != tt.want {
				t.Fatalf("ContainsWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
