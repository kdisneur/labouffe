package recipe

import (
	"fmt"
	"time"
)

// Duration is a custom duration time enabling YAML unmarshalling
type Duration time.Duration

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
