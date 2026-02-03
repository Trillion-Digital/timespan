package timespan

import "time"

type CustomWindow struct {
	start                time.Time
	end                  time.Time
	duration             time.Duration
	shouldStartBeLastDay bool
	shouldEndBeLastDay   bool
}

func (c *CustomWindow) Index() int {
	return 1
}

func (c *CustomWindow) Start() time.Time { return c.start }
func (c *CustomWindow) SetStart(t time.Time) {
	c.start = t
	c.shouldStartBeLastDay = isLastDayOfMonth(t)
}

func (c *CustomWindow) End() time.Time { return c.end }
func (c *CustomWindow) SetEnd(t time.Time) {
	c.end = t
	c.shouldEndBeLastDay = isLastDayOfMonth(t)
}

func (c *CustomWindow) Next(s ...Step) Window {
	step, ok := GetFirst(s)
	if ok {
		switch step {
		case StepYear:
			return c.shift(12)
		case StepMonth:
			return c.shift(1)
		}
	}
	return c.shiftByDuration(1)
}

func (c *CustomWindow) Prev(s ...Step) Window {
	step, ok := GetFirst(s)
	if ok {
		switch step {
		case StepYear:
			return c.shift(-12)
		case StepMonth:
			return c.shift(-1)
		}
	}
	return c.shiftByDuration(-1)
}

func (c *CustomWindow) shift(months int) Window {
	endRef := shiftMonthClamp(c.end, months)
	startRef := shiftMonthClamp(c.start, months)

	start := startRef
	end := endRef

	if c.shouldEndBeLastDay {
		end = snapToLastDayOfMonth(end)
	}
	if c.shouldStartBeLastDay {
		start = snapToLastDayOfMonth(start)
	}

	return &CustomWindow{
		start:                start,
		end:                  end,
		duration:             c.duration,
		shouldStartBeLastDay: c.shouldStartBeLastDay,
		shouldEndBeLastDay:   c.shouldEndBeLastDay,
	}
}

func (c *CustomWindow) Complete() Window {
	return c
}

func (c *CustomWindow) shiftByDuration(delta int) Window {
	d := c.duration * time.Duration(delta)

	start := c.start.Add(d)
	end := c.end.Add(d)

	return &CustomWindow{
		start:                start,
		end:                  end,
		duration:             c.duration,
		shouldStartBeLastDay: c.shouldStartBeLastDay,
		shouldEndBeLastDay:   c.shouldEndBeLastDay,
	}
}

func NewCustomWindow(start, end time.Time) Window {
	if end.Before(start) {
		panic("custom window end before start")
	}

	return &CustomWindow{
		start:                start,
		end:                  end,
		duration:             end.Sub(start),
		shouldStartBeLastDay: isLastDayOfMonth(start),
		shouldEndBeLastDay:   isLastDayOfMonth(end),
	}
}
