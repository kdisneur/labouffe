package recipe

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// IncludedIngredient represents a recipe ingredient
type IncludedIngredient struct {
	*Ingredient
	Quantity     Quantity
	Details      string
	Alternatives []*Ingredient
}

// Servings represents the number of people or number of items the recipe is made for
type Servings struct {
	Quantity int
	Type     ServingsType
}

// Recipe represents a recipe
type Recipe struct {
	Code         string
	Title        string
	Category     Category
	Cooking      Duration
	Resting      Duration
	Preparation  Duration
	Difficulty   Difficulty
	Pricing      Price
	Servings     Servings
	Instructions []*Instruction
	Ingredients  []*IncludedIngredient
	Warning      *string
}

// YAMLRecipe represents the YAML format
type YAMLRecipe struct {
	Title        string                  `yaml:"title"`
	Category     Category                `yaml:"category"`
	Cooking      Duration                `yaml:"cooking"`
	Resting      Duration                `yaml:"resting"`
	Preparation  Duration                `yaml:"preparation"`
	Difficulty   Difficulty              `yaml:"difficulty"`
	Pricing      Price                   `yaml:"pricing"`
	Servings     Servings                `yaml:"servings"`
	Instructions []*Instruction          `yaml:"instructions"`
	Ingredients  YAMLIncludedIngredients `yaml:"ingredients"`
	Warning      *string                 `yaml:"warning"`
}

// YAMLIncludedIngredient represents an ingredient to include in a recipe
type YAMLIncludedIngredient struct {
	Code         string   `yaml:"-"`
	Quantity     Quantity `yaml:"quantity"`
	Details      string   `yaml:"details"`
	Alternatives []string `yaml:"alternatives"`
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
