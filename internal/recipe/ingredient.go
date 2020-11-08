package recipe

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// Ingredient represents an actual ingredient
type Ingredient struct {
	Code  string
	Title string
}

// YAMLIngredient represents an ingredient in its YAML format
type YAMLIngredient struct {
	Code  string `yaml:"-"`
	Title string `yaml:"title"`
}

// YAMLIngredients is a type used to keep a map ordered
type YAMLIngredients []YAMLIngredient

// UnmarshalYAML is a custom unmarshaller to keep a map ordered
func (i *YAMLIngredients) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ordered := make(yaml.MapSlice, 0)
	if err := unmarshal(&ordered); err != nil {
		return fmt.Errorf("can't parse the ingredients map as yaml.MapSlice: %v", err)
	}
	unordered := make(map[string]*YAMLIngredient)
	if err := unmarshal(&unordered); err != nil {
		return fmt.Errorf("can't parse the ingredient map as map[string]YAMLIngredient: %v", err)
	}

	for _, item := range ordered {
		code, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("can't decode ingredient code as string '%v': %T", item.Key, item.Key)
		}

		unordered[code].Code = code

		*i = append(*i, *unordered[code])
	}

	return nil
}
