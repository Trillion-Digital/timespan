package timespan_test

import (
	"testing"

	"github.com/Trillion-Digital/timespan"
)

func TestWindowIndex(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "month january", got: timespan.NewMonthWindowEndingOn(mustDate(t, "2026-01-31")).Index(), want: 1},
		{name: "month december", got: timespan.NewMonthWindowEndingOn(mustDate(t, "2026-12-31")).Index(), want: 12},
		{name: "week one", got: timespan.NewWeekWindowEndingOn(mustDate(t, "2026-03-07")).Index(), want: 1},
		{name: "week four", got: timespan.NewWeekWindowEndingOn(mustDate(t, "2026-03-31")).Index(), want: 4},
		{name: "half month first half", got: timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-15")).Index(), want: 1},
		{name: "half month second half", got: timespan.NewHalfMonthWindowEndingOn(mustDate(t, "2026-03-31")).Index(), want: 2},
		{name: "quarter one", got: timespan.NewQuarterWindowEndingOn(mustDate(t, "2026-03-31")).Index(), want: 1},
		{name: "quarter four", got: timespan.NewQuarterWindowEndingOn(mustDate(t, "2026-12-31")).Index(), want: 4},
		{name: "semester one", got: timespan.NewSemesterWindowEndingOn(mustDate(t, "2026-06-30")).Index(), want: 1},
		{name: "semester two", got: timespan.NewSemesterWindowEndingOn(mustDate(t, "2026-12-31")).Index(), want: 2},
		{name: "year is always one", got: timespan.NewYearWindowEndingOn(mustDate(t, "2026-12-31")).Index(), want: 1},
		{name: "custom is always one", got: timespan.NewCustomWindow(mustDate(t, "2026-01-10"), mustDate(t, "2026-01-25")).Index(), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("Index() = %d, want %d", tt.got, tt.want)
			}
		})
	}
}
