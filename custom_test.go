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
		{
			name:      "simple custom window",
			start:     "2026-01-10",
			end:       "2026-01-20",
			wantStart: "2026-01-10",
			wantEnd:   "2026-01-20",
		},
		{
			name:      "custom window ending on last day",
			start:     "2026-01-05",
			end:       "2026-01-31",
			wantStart: "2026-01-05",
			wantEnd:   "2026-01-31",
		},
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

func TestCustomWindow_Next_Month(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "same days next month",
			start:     "2026-01-10",
			end:       "2026-01-20",
			wantStart: "2026-02-10",
			wantEnd:   "2026-02-20",
		},
		{
			name:      "end clamped to february",
			start:     "2026-01-05",
			end:       "2026-01-31",
			wantStart: "2026-02-05",
			wantEnd:   "2026-02-28",
		},
		{
			name:      "start clamped independently",
			start:     "2026-01-31",
			end:       "2026-02-15",
			wantStart: "2026-02-28",
			wantEnd:   "2026-03-15",
		},
		{
			name:      "both start and end last day preserved",
			start:     "2026-01-31",
			end:       "2026-02-28",
			wantStart: "2026-02-28",
			wantEnd:   "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(
				mustDate(t, tt.start),
				mustDate(t, tt.end),
			)

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

func TestCustomWindow_Prev_Month(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "same days prev month",
			start:     "2026-03-10",
			end:       "2026-03-20",
			wantStart: "2026-02-10",
			wantEnd:   "2026-02-20",
		},
		{
			name:      "end last day preserved backwards",
			start:     "2026-03-05",
			end:       "2026-03-31",
			wantStart: "2026-02-05",
			wantEnd:   "2026-02-28",
		},
		{
			name:      "both start and end last day backwards",
			start:     "2026-03-31",
			end:       "2026-04-30",
			wantStart: "2026-02-28",
			wantEnd:   "2026-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(
				mustDate(t, tt.start),
				mustDate(t, tt.end),
			)

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
func TestCustomWindow_Next_Year(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "simple year shift",
			start:     "2026-02-10",
			end:       "2026-03-05",
			wantStart: "2027-02-10",
			wantEnd:   "2027-03-05",
		},
		{
			name:      "leap year clamped",
			start:     "2024-02-29",
			end:       "2024-03-31",
			wantStart: "2025-02-28",
			wantEnd:   "2025-03-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(
				mustDate(t, tt.start),
				mustDate(t, tt.end),
			)

			got := w.Next(timespan.StepYear)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestCustomWindow_Prev_Year(t *testing.T) {
	tests := []struct {
		name      string
		start     string
		end       string
		wantStart string
		wantEnd   string
	}{
		{
			name:      "simple year back",
			start:     "2026-02-10",
			end:       "2026-03-05",
			wantStart: "2025-02-10",
			wantEnd:   "2025-03-05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := timespan.NewCustomWindow(
				mustDate(t, tt.start),
				mustDate(t, tt.end),
			)

			got := w.Prev(timespan.StepYear)

			assertWindow(
				t,
				got,
				mustDate(t, tt.wantStart),
				mustDate(t, tt.wantEnd),
			)
		})
	}
}

func TestCustomWindow_Complete(t *testing.T) {
	w := timespan.NewCustomWindow(
		mustDate(t, "2026-01-10"),
		mustDate(t, "2026-01-25"),
	)

	got := w.Complete()

	assertWindow(
		t,
		got,
		mustDate(t, "2026-01-10"),
		mustDate(t, "2026-01-25"),
	)
}

func TestCustomWindow_DurationInvariant(t *testing.T) {
	start := mustDate(t, "2026-01-05")
	end := mustDate(t, "2026-01-20")

	w := timespan.NewCustomWindow(start, end)

	n := w.Next()
	p := w.Prev()

	if n.End().Sub(n.Start()) != end.Sub(start) {
		t.Fatalf("duration changed on Next")
	}
	if p.End().Sub(p.Start()) != end.Sub(start) {
		t.Fatalf("duration changed on Prev")
	}
}
