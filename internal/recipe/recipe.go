package recipe

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// IncludedIngredient represents a recipe ingredient
type IncludedIngredient struct {
	Ingredient
	Quantity Quantity
}

// Recipe represents a recipe
type Recipe struct {
	Code         string
	Title        string
	Cooking      Duration
	Preparation  Duration
	Difficulty   Difficulty
	Pricing      Price
	Guests       int
	Instructions []string
	Ingredients  []IncludedIngredient
}

// YAMLRecipe represents the YAML format
type YAMLRecipe struct {
	Title        string                  `yaml:"title"`
	Cooking      Duration                `yaml:"cooking"`
	Preparation  Duration                `yaml:"preparation"`
	Difficulty   Difficulty              `yaml:"difficulty"`
	Pricing      Price                   `yaml:"pricing"`
	Guests       int                     `yaml:"guests"`
	Instructions []string                `yaml:"instructions"`
	Ingredients  YAMLIncludedIngredients `yaml:"ingredients"`
}

// YAMLIncludedIngredient represents an ingredient to include in a recipe
type YAMLIncludedIngredient struct {
	Code     string   `yaml:"-"`
	Quantity Quantity `yaml:"quantity"`
}

// YAMLIncludedIngredients is a type used to keep a map ordered
type YAMLIncludedIngredients []YAMLIncludedIngredient

// UnmarshalYAML is a custom unmarshaller to keep a map ordered
func (i *YAMLIncludedIngredients) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ordered := make(yaml.MapSlice, 0)
	if err := unmarshal(&ordered); err != nil {
		return fmt.Errorf("can't parse the included ingredient map as yaml.MapSlice: %v", err)
	}
	unordered := make(map[string]*YAMLIncludedIngredient)
	if err := unmarshal(&unordered); err != nil {
		return fmt.Errorf("can't parse the included ingredient map as map[string]includedingredient: %v", err)
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
