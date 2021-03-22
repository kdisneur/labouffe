package recipe

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Builder represents a recipe builder. It can parse and add new recipes
// ad long as the ingredients already exist
type Builder struct {
	Ingredients []*Ingredient
	Recipes     []*Recipe
}

// NewBuilderFromYAMLIngredients parses a YAML list of ingredients
func NewBuilderFromYAMLIngredients(input io.Reader) (*Builder, error) {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	var b Builder

	var ingredients YAMLIngredients
	if err := decoder.Decode(&ingredients); err != nil {
		return nil, fmt.Errorf("can't decode ingredients: %v", err)
	}

	for _, data := range ingredients {
		b.Ingredients = append(b.Ingredients, &Ingredient{
			Code:       data.Code,
			Title:      data.Title,
			recipeCode: data.RecipeCode,
		})
	}

	return &b, nil
}

// ParseNewYAMLRecipe parses a new recipe from a YAML
func (b *Builder) ParseNewYAMLRecipe(code string, input io.Reader) error {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	var r YAMLRecipe
	if err := decoder.Decode(&r); err != nil {
		return fmt.Errorf("can't decode recipe: %v", err)
	}

	ingredients := make([]*IncludedIngredient, len(r.Ingredients))
	for i, item := range r.Ingredients {
		ingredient, err := b.findIngredientByCode(item.Code)
		if err != nil {
			return fmt.Errorf("can't find ingredient '%s'", item.Code)
		}

		alternatives := make([]*Ingredient, len(item.Alternatives))
		for i := range item.Alternatives {
			ingredient, err := b.findIngredientByCode(item.Alternatives[i])
			if err != nil {
				return fmt.Errorf("can't find ingredient '%s' as alternative to '%s'", item.Alternatives[i], item.Code)
			}

			alternatives[i] = ingredient
		}

		ingredients[i] = &IncludedIngredient{
			Ingredient:   ingredient,
			Quantity:     item.Quantity,
			Details:      item.Details,
			Alternatives: alternatives,
		}
	}

	b.Recipes = append(b.Recipes, &Recipe{
		Code:         code,
		Category:     r.Category,
		Title:        r.Title,
		Cooking:      r.Cooking,
		Preparation:  r.Preparation,
		Difficulty:   r.Difficulty,
		Pricing:      r.Pricing,
		Servings:     r.Servings,
		Instructions: r.Instructions,
		Ingredients:  ingredients,
		Warning:      r.Warning,
	})

	return nil
}

// LoadRecipeIngredients is a function that associate a recipe with
// an ingredient. This function must be called only after the recipes
// have already been loaded
func (b *Builder) LoadRecipeIngredients() error {
	for i := range b.Ingredients {
		ingredient := b.Ingredients[i]
		if ingredient.recipeCode == "" {
			continue
		}

		recipe, err := b.findRecipeByCode(ingredient.recipeCode)
		if err != nil {
			return fmt.Errorf("can't associate ingredient '%s' with recipe '%s': %v", ingredient.Code, ingredient.recipeCode, err)
		}

		ingredient.Recipe = recipe
		if ingredient.Title == "" {
			ingredient.Title = recipe.Title
		}
	}

	return nil
}

// LoadRecipeInstructions is a function that associate a recipe with
// a recipe step. This function must be called only after the recipes
// have already been loaded
func (b *Builder) LoadRecipeInstructions() error {
	for i := range b.Recipes {
		recipe := b.Recipes[i]
		for i := range recipe.Instructions {
			instruction := recipe.Instructions[i]
			if instruction.recipeCode == "" {
				continue
			}

			associatedRecipe, err := b.findRecipeByCode(instruction.recipeCode)
			if err != nil {
				return fmt.Errorf("recipe '%s': can't associate recipe '%s' for step '%s': %v", recipe.Code, instruction.recipeCode, instruction.Title, err)
			}

			instruction.Recipe = associatedRecipe
		}
	}

	return nil
}

// Length returns the number of parsed recipes
func (b *Builder) Length() int {
	return len(b.Recipes)
}

func (b *Builder) findIngredientByCode(code string) (*Ingredient, error) {
	for _, ingredient := range b.Ingredients {
		if ingredient.Code == code {
			return ingredient, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (b *Builder) findRecipeByCode(code string) (*Recipe, error) {
	for _, recipe := range b.Recipes {
		if recipe.Code == code {
			return recipe, nil
		}
	}

	return nil, fmt.Errorf("not found")
}
