package recipe

import (
	"fmt"
	"strings"
)

// Difficulty represents how hard to make is a recipe
type Difficulty int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=Difficulty -linecomment
const (
	// DifficultyEasy represents a recipe that is easy to make
	DifficultyEasy Difficulty = iota // easy
	// DifficultyAverage represents a recipe that is not too hard to make
	DifficultyAverage // average
	// DifficultyHard represents a recipe that is hard to make
	DifficultyHard // hard
)

// UnmarshalYAML is the function in charge of unmarshalling the string value to a Go constant
func (d *Difficulty) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var namedDifficulty string
	if err := unmarshal(&namedDifficulty); err != nil {
		return fmt.Errorf("can't unmarshal named difficulty, expected a string: %v", err)
	}

	var allNamedDifficulty []string
	for i := 0; i < len(_Difficulty_index)-1; i++ {
		allNamedDifficulty = append(allNamedDifficulty, _Difficulty_name[_Difficulty_index[i]:_Difficulty_index[i+1]])
	}

	for i := range allNamedDifficulty {
		if allNamedDifficulty[i] == strings.ToLower(namedDifficulty) {
			*d = Difficulty(i)
			return nil
		}
	}

	return fmt.Errorf("unsupported difficulty '%s'. valid difficulty are: %s", namedDifficulty, strings.Join(allNamedDifficulty, ", "))
}
