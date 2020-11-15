package recipe

import (
	"fmt"
	"math"
	"time"
)

// DurationRange represents a duration group
type DurationRange struct {
	ThresholdMinutes float64
	Title            string
}

// AllDurationRanges represents all the available duration ranges
func AllDurationRanges() []DurationRange {
	return []DurationRange{
		{ThresholdMinutes: 30, Title: "< 30m"},
		{ThresholdMinutes: 60, Title: "> 30m et < 1h"},
		{ThresholdMinutes: 90, Title: "> 1h et < 1h30"},
		{ThresholdMinutes: math.MaxFloat64, Title: "> 1h30"},
	}
}

// Duration is a custom duration time enabling YAML unmarshalling
type Duration time.Duration

// TimeDuration returns the underlying time.Duration
func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
}

// RangeName groups the duration in useful filters
func (d Duration) RangeName() string {
	allRanges := AllDurationRanges()
	for _, r := range allRanges {
		if r.ThresholdMinutes > time.Duration(d).Minutes() {
			return r.Title
		}
	}

	return allRanges[len(allRanges)-1].Title
}

func (d Duration) String() string {
	return time.Duration(d).String()
}

// UnmarshalYAML is the function in charge of unmarshalling the string value to a Duration
func (d *Duration) UnmarshalYAML(unmarshal func(input interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return fmt.Errorf("can't read timestamp as string: %v", err)
	}

	t, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("can't parse '%s' to Duration: %v", t, err)
	}

	*d = Duration(t)

	return nil
}
