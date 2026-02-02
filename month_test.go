package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

func TestMonthlyWindow_Next(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "end mid month moves normally",
			input:     mustDate(t, "2026-01-15"),
			wantStart: "2026-02-01",
			wantEnd:   "2026-02-15",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "end on 31 clamps to feb and snaps back",
			input:     mustDate(t, "2026-01-31"),
			wantStart: "2026-02-01",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "feb last day snaps back to march 31",
			input:     mustDate(t, "2026-02-28"),
			wantStart: "2026-03-01",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "feb mid day goes to to mid march on last day",
			input:     mustDate(t, "2026-02-15"),
			wantStart: "2026-03-15",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewMonthWindowStartingOn,
		},
		{
			name:      "feb last day goes march to jan on last day",
			input:     mustDate(t, "2026-02-28"),
			wantStart: "2026-03-31",
			wantEnd:   "2026-03-31",
			fn:        timespan.NewMonthWindowStartingOn,
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

func TestMonthlyWindow_Prev(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
		fn        func(time.Time) timespan.Window
	}{
		{
			name:      "end mid month moves back normally",
			input:     mustDate(t, "2026-03-15"),
			wantStart: "2026-02-01",
			wantEnd:   "2026-02-15",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "end on 31 goes back to feb last day",
			input:     mustDate(t, "2026-03-31"),
			wantStart: "2026-02-01",
			wantEnd:   "2026-02-28",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "feb last day goes back to jan 31",
			input:     mustDate(t, "2026-02-28"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-01-31",
			fn:        timespan.NewMonthWindowEndingOn,
		},
		{
			name:      "feb mid day goes back to mid jan on last day",
			input:     mustDate(t, "2026-02-15"),
			wantStart: "2026-01-15",
			wantEnd:   "2026-01-31",
			fn:        timespan.NewMonthWindowStartingOn,
		},
		{
			name:      "feb last day goes back to jan on last day",
			input:     mustDate(t, "2026-02-28"),
			wantStart: "2026-01-31",
			wantEnd:   "2026-01-31",
			fn:        timespan.NewMonthWindowStartingOn,
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

func TestMonthlyWindow_Complete(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart string
		wantEnd   string
	}{
		{
			name:      "mid month completes to full month",
			input:     mustDate(t, "2026-01-15"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-01-31",
		},
		{
			name:      "already full month remains full",
			input:     mustDate(t, "2026-01-31"),
			wantStart: "2026-01-01",
			wantEnd:   "2026-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewMonthWindowEndingOn(tt.input)
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

func TestNewMonthlyWindowStartingOn(t *testing.T) {
	t.Run("starting mid month", func(t *testing.T) {
		got := timespan.NewMonthWindowStartingOn(mustDate(t, "2026-01-10"))

		assertWindow(
			t,
			got,
			mustDate(t, "2026-01-10"),
			mustDate(t, "2026-01-31"),
		)
	})
}

func TestNewMonthlyWindowEndingOn(t *testing.T) {
	t.Run("ending mid month", func(t *testing.T) {
		got := timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-10"))

		assertWindow(
			t,
			got,
			mustDate(t, "2026-01-01"),
			mustDate(t, "2026-01-10"),
		)
	})

	t.Run("ending on last day", func(t *testing.T) {
		got := timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-31"))

		assertWindow(
			t,
			got,
			mustDate(t, "2026-01-01"),
			mustDate(t, "2026-01-31"),
		)
	})
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("invalid date %q: %v", s, err)
	}
	return d
}

func assertWindow(t *testing.T, got timespan.Window, wantStart, wantEnd time.Time) {
	t.Helper()

	if !got.Start().Equal(wantStart) {
		t.Errorf("start = %v, want %v", got.Start(), wantStart)
	}
	if !got.End().Equal(wantEnd) {
		t.Errorf("end = %v, want %v", got.End(), wantEnd)
	}
}
