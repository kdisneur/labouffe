package recipe

import (
	"fmt"
)

// Instruction represents a step of a recipe
type Instruction struct {
	Title      string
	Recipe     *Recipe
	Warning    string
	recipeCode string
}

// UnmarshalYAML unmarshals a string or an object depending of the format
func (i *Instruction) UnmarshalYAML(unmarhsal func(interface{}) error) error {
	var title string
	var recipeCode string
	var warning string

	if err1 := unmarhsal(&title); err1 != nil {
		var data map[string]string
		if err2 := unmarhsal(&data); err2 != nil {
			return fmt.Errorf("can't load instruction: %v OR %v", err1, err2)
		}

		title = data["instruction"]
		recipeCode = data["recipe"]
		warning = data["warning"]
	}

	i.Title = title
	i.recipeCode = recipeCode
	i.Warning = warning

	return nil
}
