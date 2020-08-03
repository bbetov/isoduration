// Package isoduration implements ISO8601 durations.

// See https://en.wikipedia.org/wiki/ISO_8601#Durations

package isoduration

import (
	"reflect"
	"testing"
)

func TestParseWeeks(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    *Duration
		wantErr bool
	}{
		{"P4W", "P4W", &Duration{Milliseconds: 2419200 * millisecondsPerSecond}, false},
		{"-P4W", "-P4W", &Duration{Milliseconds: -2419200 * millisecondsPerSecond}, false},
		{"+P-4W", "+P-4W", &Duration{Milliseconds: -2419200 * millisecondsPerSecond}, false},
		{"-P-4W", "-P-4W", &Duration{Milliseconds: 2419200 * millisecondsPerSecond}, false},
		{"P-4W", "P-4W", &Duration{Milliseconds: -2419200 * millisecondsPerSecond}, false},
		{"P-4WT54S", "P-4WT54S", nil, true},
		{"P4WT13M", "P4WT13M", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFull(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    *Duration
		wantErr bool
	}{
		{"P23DT23H", "P23DT23H", &Duration{Milliseconds: millisecondsPerSecond * (23*secondsPerDay + 23*secondsPerHour)}, false},
		{"P4Y", "P4Y", &Duration{Milliseconds: millisecondsPerSecond * (4 * secondsPerYear)}, false},
		{"P1DT12H", "P1DT12H", &Duration{Milliseconds: millisecondsPerSecond * (secondsPerDay + 12*secondsPerHour)}, false},
		{"PT36H", "PT36H", &Duration{Milliseconds: millisecondsPerSecond * (36 * secondsPerHour)}, false},
		{"PT36H34M12.987S", "PT36H34M12.987S", &Duration{Milliseconds: millisecondsPerSecond*(36*secondsPerHour+34*secondsPerMinute+12) + 987}, false},
		{"PT36H-34M-12.987S", "PT36H-34M-12.987S", &Duration{Milliseconds: millisecondsPerSecond*(36*secondsPerHour-34*secondsPerMinute-12) - 987}, false},
		{"PT36.0H-34M-12.987S", "PT36.0H-34M-12.987S", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
