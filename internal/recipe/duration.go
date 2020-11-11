package recipe

import (
	"fmt"
	"time"
)

// Duration is a custom duration time enabling YAML unmarshalling
type Duration time.Duration

// TimeDuration returns the underlying time.Duration
func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
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
