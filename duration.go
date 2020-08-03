// Package isoduration implements ISO8601 durations.
// See https://en.wikipedia.org/wiki/ISO_8601#Durations
package isoduration

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Duration holds the data needed to hold duration
type Duration struct {
	Milliseconds int64
}

var (
	errorInvalidDuration = errors.New("Invalid duration string")
)

// Parse converts an ISO8601 duration string to a Duration object
func Parse(duration string) (*Duration, error) {
	duration = strings.ToUpper(duration)

	if weeks.MatchString(duration) {
		return parseWeeks(duration)
	} else if full.MatchString(duration) {
		return parseFull(duration)
	}

	return nil, errorInvalidDuration
}

// String returns a string representation of a duration, where the lagrest time interval is days.
// This is due to the fact that the larger time intervals have lengths, dependent on other factors.
// The individual components will always have positive sign, only the whole duration's sign changes.
func (d *Duration) String() string {

	if d.Milliseconds == 0 {
		return "PT0S"
	}
	sign := ""
	if d.Milliseconds < 0 {
		sign = "-"
		d.Milliseconds = -d.Milliseconds
	}

	var b strings.Builder
	b.Grow(64)
	fmt.Fprintf(&b, "%sP", sign)

	secs := d.Milliseconds / millisecondsPerSecond
	mills := d.Milliseconds % millisecondsPerSecond

	if secs > secondsPerDay {
		days := secs / secondsPerDay
		fmt.Fprintf(&b, "%dD", days)
		secs -= days * secondsPerDay
	}

	if secs > 0 || mills > 0 {
		for secs > 0 || mills > 0 {
			switch {
			case secs > secondsPerHour:
				hours := secs / secondsPerHour
				fmt.Fprintf(&b, "%dH", hours)
				secs -= hours * secondsPerHour
			case secs > secondsPerMinute:
				mins := secs / secondsPerMinute
				fmt.Fprintf(&b, "%dH", mins)
				secs -= mins * secondsPerMinute
			default:
				fmt.Fprintf(&b, "%d.%dS", secs, mills)
				break
			}
		}
	}
	return b.String()
}

// StringWeeks returns a truncated number of weeks (i.e. for 1.5 weeks, it will return 1 week)
func (d *Duration) StringWeeks() string {
	sign := ""
	if d.Milliseconds < 0 {
		sign = "-"
		d.Milliseconds = -d.Milliseconds
	}

	if d.Milliseconds < secondsPerWeek*millisecondsPerSecond {
		return "P0W"
	}
	return fmt.Sprintf("%sP%dW", sign, (d.Milliseconds/millisecondsPerSecond)/secondsPerWeek)
}

// Duration converts isoduration.Duration to time.Duration
// Since time.Duration is in nanonseconds and isoduration.Duration is microseconds
// it is possible that the time.Duration overflows
func (d *Duration) Duration() (time.Duration, error) {
	if d.Milliseconds >= math.MaxInt64/1000 || d.Milliseconds <= math.MinInt64/1000 {
		return time.Duration(0), errors.New("Overflow converting isoduration.Duration to time.Duration")
	}

	return time.Duration(d.Milliseconds * 1000), nil
}

// FromDuration creates a duration object from time.Duration, truncating the nanoseconds
func FromDuration(d time.Duration) *Duration {
	return &Duration{Milliseconds: d.Milliseconds()}
}
