package foodaccess

import (
	"fmt"
	"strings"
)

// ServingsType represents how hard to make is a recipe
type ServingsType int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=ServingsType -linecomment
const (
	// ServingsTypeGuests represents a serving counted by number of guests
	ServingsTypeGuests ServingsType = iota // personnes
	// ServingsTypeUnits represents a serving conted by number of units
	ServingsTypeUnits // unités
)

// AllServingsType returns the full list of servingstypes
func AllServingsType() []ServingsType {
	var servingstypes []ServingsType
	for i := 0; i < len(_ServingsType_index)-1; i++ {
		servingstypes = append(servingstypes, ServingsType(i))
	}

	return servingstypes
}

// UnmarshalYAML is the function in charge of unmarshalling the string value to a Go constant
func (d *ServingsType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var namedServingsType string
	if err := unmarshal(&namedServingsType); err != nil {
		return fmt.Errorf("can't unmarshal named difficulty, expected a string: %v", err)
	}

	var allNamedServingsType []string
	for i := 0; i < len(_ServingsType_index)-1; i++ {
		allNamedServingsType = append(allNamedServingsType, _ServingsType_name[_ServingsType_index[i]:_ServingsType_index[i+1]])
	}

	for i := range allNamedServingsType {
		if allNamedServingsType[i] == strings.ToLower(namedServingsType) {
			*d = ServingsType(i)
			return nil
		}
	}

	return fmt.Errorf("unsupported difficulty '%s'. valid difficulty are: %s", namedServingsType, strings.Join(allNamedServingsType, ", "))
}
