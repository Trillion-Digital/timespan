package timespan_test

import (
	"testing"

	"github.com/Trillion-Digital/timespan"
)

func TestNewCustomWindow(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		wantStart string
		wantEnd   string
	}{
		{name: "simple custom window", start: "2026-01-10", end: "2026-01-20", wantStart: "2026-01-10", wantEnd: "2026-01-20"},
		{name: "custom window ending on last day", start: "2026-01-05", end: "2026-01-31", wantStart: "2026-01-05", wantEnd: "2026-01-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(
				mustDate(t, tt.start),
				mustDate(t, tt.end),
			)

			assertWindow(
				t,
				w,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestNewCustomWindow_PanicsWhenEndBeforeStart(t *testing.T) {
	assertPanics(t, func() {
		timespan.NewCustomWindow(mustDate(t, "2026-01-20"), mustDate(t, "2026-01-10"))
	})
}

func TestCustomWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		step      timespan.Step
		wantStart string
		wantEnd   string
	}{
		{name: "default shifts by window duration", start: "2026-01-10", end: "2026-01-20", wantStart: "2026-01-20", wantEnd: "2026-01-30"},
		{name: "step month keeps same calendar days", start: "2026-01-10", end: "2026-01-20", step: timespan.StepMonth, wantStart: "2026-02-10", wantEnd: "2026-02-20"},
		{name: "step month clamps and preserves end of month", start: "2026-01-05", end: "2026-01-31", step: timespan.StepMonth, wantStart: "2026-02-05", wantEnd: "2026-02-28"},
		{name: "step year clamps leap day", start: "2024-02-29", end: "2024-03-31", step: timespan.StepYear, wantStart: "2025-02-28", wantEnd: "2025-03-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(mustDate(t, tt.start), mustDate(t, tt.end))

			var got timespan.Window
			if tt.step.Valid() {
				got = w.Next(tt.step)
			} else {
				got = w.Next()
			}

			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
		})
	}
}

func TestCustomWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		step      timespan.Step
		wantStart string
		wantEnd   string
	}{
		{name: "default shifts by window duration", start: "2026-03-10", end: "2026-03-20", wantStart: "2026-02-28", wantEnd: "2026-03-10"},
		{name: "step month keeps same calendar days", start: "2026-03-10", end: "2026-03-20", step: timespan.StepMonth, wantStart: "2026-02-10", wantEnd: "2026-02-20"},
		{name: "step month preserves last days", start: "2026-03-31", end: "2026-04-30", step: timespan.StepMonth, wantStart: "2026-02-28", wantEnd: "2026-03-31"},
		{name: "step year moves back one calendar year", start: "2026-02-10", end: "2026-03-05", step: timespan.StepYear, wantStart: "2025-02-10", wantEnd: "2025-03-05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(mustDate(t, tt.start), mustDate(t, tt.end))

			var got timespan.Window
			if tt.step.Valid() {
				got = w.Prev(tt.step)
			} else {
				got = w.Prev()
			}

			assertWindow(t, got, mustDate(t, tt.wantStart), mustDate(t, tt.wantEnd))
		})
	}
}

func TestCustomWindow_Complete(t *testing.T) {
	w := timespan.NewCustomWindow(mustDate(t, "2026-01-10"), mustDate(t, "2026-01-25"))
	got := w.Complete()
	assertWindow(t, got, mustDate(t, "2026-01-10"), mustDate(t, "2026-01-25"))
}

func TestCustomWindow_DurationInvariant(t *testing.T) {
	start := mustDate(t, "2026-01-05")
	end := mustDate(t, "2026-01-20")
	w := timespan.NewCustomWindow(start, end)

	if got := w.Next().End().Sub(w.Next().Start()); got != end.Sub(start) {
		t.Fatalf("Next duration = %v, want %v", got, end.Sub(start))
	}
	if got := w.Prev().End().Sub(w.Prev().Start()); got != end.Sub(start) {
		t.Fatalf("Prev duration = %v, want %v", got, end.Sub(start))
	}
}
