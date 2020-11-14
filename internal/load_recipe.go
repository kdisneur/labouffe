package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kdisneur/labouffe/internal/recipe"
)

// LoadIngredientAndRecipes is in charge of loading / parsing / validating all ingredients and recipes data
func LoadIngredientAndRecipes(ingredientpath string, recipesfolderpath string) ([]*recipe.Ingredient, []*recipe.Recipe, error) {
	ingredientfile, err := os.Open(ingredientpath)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open ingredient file '%s': %v", ingredientpath, err)
	}
	defer ingredientfile.Close()

	builder, err := recipe.NewBuilderFromYAMLIngredients(ingredientfile)
	if err != nil {
		return nil, nil, fmt.Errorf("can't import ingredients file '%s': %v", ingredientpath, err)
	}

	err = filepath.Walk(recipesfolderpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		basename := filepath.Base(path)
		extension := filepath.Ext(path)

		if extension != ".yaml" && extension != ".yml" {
			return nil
		}

		recipefile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("can't open recipe file: '%s': %v", path, err)
		}
		defer recipefile.Close()

		code := strings.TrimSuffix(basename, extension)
		if err := builder.ParseNewYAMLRecipe(code, recipefile); err != nil {
			return fmt.Errorf("can't import recipe '%s': %v", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("can't import recipes from '%s': %v", recipesfolderpath, err)
	}

	if err := builder.LoadRecipeIngredients(); err != nil {
		return nil, nil, fmt.Errorf("can't validate ingredients: %v", err)
	}

	if err := builder.LoadRecipeInstructions(); err != nil {
		return nil, nil, fmt.Errorf("can't validate recipe instructions: %v", err)
	}

	return builder.Ingredients, builder.Recipes, nil
}
