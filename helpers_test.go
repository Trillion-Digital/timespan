package timespan_test

import (
	"testing"
	"time"

	"github.com/Trillion-Digital/timespan"
)

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

func collectDays(seq func(func(time.Time) bool)) []time.Time {
	var days []time.Time
	seq(func(day time.Time) bool {
		days = append(days, day)
		return true
	})
	return days
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}
