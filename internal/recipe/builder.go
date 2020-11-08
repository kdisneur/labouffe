package recipe

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// Builder represents a recipe builder. It can parse and add new recipes
// ad long as the ingredients already exist
type Builder struct {
	Ingredients []Ingredient
	Recipes     []Recipe
}

// NewBuilderFromYAMLIngredients creates a new builder for ingredients and recipes loading
// a YAML list of ingredients
func NewBuilderFromYAMLIngredients(input io.Reader) (*Builder, error) {
	var b Builder

	var ingredients YAMLIngredients
	if err := yaml.NewDecoder(input).Decode(&ingredients); err != nil {
		return nil, fmt.Errorf("can't decode ingredients: %v", err)
	}

	for _, data := range ingredients {
		b.Ingredients = append(b.Ingredients, Ingredient(data))
	}

	return &b, nil
}

// ParseNewYAMLRecipe parses a new recipe from a YAML
func (b *Builder) ParseNewYAMLRecipe(code string, input io.Reader) error {
	var r YAMLRecipe
	if err := yaml.NewDecoder(input).Decode(&r); err != nil {
		return fmt.Errorf("can't decode recipe: %v", err)
	}

	ingredients := make([]IncludedIngredient, len(r.Ingredients))
	for i, item := range r.Ingredients {
		ingredient, err := b.findIngredientByCode(item.Code)
		if err != nil {
			return fmt.Errorf("can't find ingredient '%s'", item.Code)
		}

		ingredients[i] = IncludedIngredient{
			Ingredient: ingredient,
			Quantity:   item.Quantity,
		}
	}

	b.Recipes = append(b.Recipes, Recipe{
		Code:         code,
		Title:        r.Title,
		Cooking:      r.Cooking,
		Preparation:  r.Preparation,
		Difficulty:   r.Difficulty,
		Pricing:      r.Pricing,
		Guests:       r.Guests,
		Instructions: r.Instructions,
		Ingredients:  ingredients,
	})

	return nil
}

// Length returns the number of parsed recipes
func (b *Builder) Length() int {
	return len(b.Recipes)
}

func (b *Builder) findIngredientByCode(code string) (Ingredient, error) {
	for _, ingredient := range b.Ingredients {
		if ingredient.Code == code {
			return ingredient, nil
		}
	}

	return Ingredient{}, fmt.Errorf("not found")
}
