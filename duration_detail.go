package isoduration

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	secondsPerMinute int64 = 60
	secondsPerHour         = 60 * secondsPerMinute
	secondsPerDay          = 24 * secondsPerHour
	secondsPerWeek         = 7 * secondsPerDay
	secondsPerYear         = 365 * secondsPerDay
	secondsPerMonth        = 30 * secondsPerDay

	millisecondsPerSecond int64 = 1000
)

var (
	// Weeks formatting (P[n]W). Counted as 7 days, 24h each day
	weeks = regexp.MustCompile(`^(?P<mainsign>[+-]?)P((?P<weeks>[+-]?\d+)W)$`)

	// full formatting
	full = regexp.MustCompile(`^(?P<mainsign>[+-]?)P((?P<years>[+-]?\d+)Y)?((?P<months>[+-]?\d+)M)?((?P<days>[+-]?\d+)D)?(T((?P<hours>[+-]?\d+)H)?((?P<minutes>[+-]?\d+)M)?((?P<secondssign>[+-]?)(?P<seconds>\d+)([,.](?P<milliseconds>\d+))?S)?)?$`)
)

func parseWeeks(duration string) (*Duration, error) {
	matches := weeks.FindStringSubmatch(duration)
	signer := int64(1)
	for i, name := range weeks.SubexpNames() {
		if i == 0 || len(name) == 0 || len(matches[i]) == 0 {
			continue
		}
		if name == "mainsign" {
			if matches[i] == "-" {
				signer = -1
			}
		} else if name == "weeks" {
			val, err := strconv.ParseInt(matches[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Unable to convert weeks value %s to integer: %v", matches[i], err)
			}
			return &Duration{Milliseconds: millisecondsPerSecond * secondsPerWeek * signer * val}, nil
		}
	}
	return nil, fmt.Errorf("Invalid weeks duration: %s", duration)
}

func parseFull(duration string) (*Duration, error) {
	matches := full.FindStringSubmatch(duration)
	signer := int64(1)
	secondssigner := int64(1)
	seconds := int64(0)
	secondsonly := int64(0)
	milliseconds := int64(0)
	for i, name := range full.SubexpNames() {
		if i == 0 || len(name) == 0 || len(matches[i]) == 0 {
			continue
		}

		if name == "mainsign" {
			if matches[i] == "-" {
				signer = -1
			}
			continue
		} else if name == "secondssign" {
			if matches[i] == "-" {
				secondssigner = -1
			}
			continue
		}
		// Assume everything else is numbers

		val, err := strconv.ParseInt(matches[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Unable to convert %s value %s to integer: %v", name, matches[i], err)
		}

		switch name {
		case "years":
			seconds += val * secondsPerYear
		case "months":
			seconds += val * secondsPerMonth
		case "days":
			seconds += val * secondsPerDay
		case "hours":
			seconds += val * secondsPerHour
		case "minutes":
			seconds += val * secondsPerMinute
		case "seconds":
			secondsonly = val
		case "milliseconds":
			if val >= millisecondsPerSecond {
				return nil, errors.New("Number of milliseconds may only be up to 1000")
			}
			milliseconds = val
		}
	}

	return &Duration{Milliseconds: signer * (millisecondsPerSecond*(seconds+secondssigner*secondsonly) + secondssigner*milliseconds)}, nil
}
