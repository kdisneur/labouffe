package recipe_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/kdisneur/labouffe/internal/recipe"
)

func TestBuilderParsing(t *testing.T) {
	type RecipeInput struct {
		Code string
		Path string
	}

	tcs := []struct {
		Name             string
		IngredientsInput string
		RecipeInputs     []RecipeInput
		ExpectedOutput   []recipe.Recipe
	}{
		{
			Name:             "with a list of valid attributes",
			IngredientsInput: "testdata/working-ingredients.yaml",
			RecipeInputs: []RecipeInput{
				{Code: "sugar-pasta", Path: "testdata/sugar-pasta.yaml"},
				{Code: "pesto-pasta", Path: "testdata/pesto-pasta.yaml"},
			},
			ExpectedOutput: []recipe.Recipe{
				{
					Code:        "sugar-pasta",
					Title:       "Sugar pasta",
					Preparation: recipe.Duration(time.Duration(2*time.Hour + 30*time.Minute)),
					Cooking:     recipe.Duration(time.Duration(10 * time.Minute)),
					Difficulty:  recipe.DifficultyEasy,
					Pricing:     recipe.PriceCheap,
					Guests:      3,
					Ingredients: []recipe.IncludedIngredient{
						{
							Ingredient: recipe.Ingredient{Code: "pasta", Title: "Pasta"},
							Quantity:   recipe.Quantity{Value: 50, Unit: recipe.QuantityUnitGram},
						},
						{
							Ingredient: recipe.Ingredient{Code: "sugar", Title: "Sugar"},
							Quantity:   recipe.Quantity{Value: 20, Unit: recipe.QuantityUnitGram},
						},
					},
					Instructions: []string{
						"Boil the pasta, without salt",
						"Drain the pasta",
						"Combine brown sugar with pasta",
					},
				},
				{
					Code:        "pesto-pasta",
					Title:       "Pesto pasta",
					Preparation: recipe.Duration(time.Duration(5 * time.Minute)),
					Cooking:     recipe.Duration(time.Duration(10 * time.Minute)),
					Difficulty:  recipe.DifficultyEasy,
					Pricing:     recipe.PriceCheap,
					Guests:      3,
					Ingredients: []recipe.IncludedIngredient{
						{
							Ingredient: recipe.Ingredient{Code: "pasta", Title: "Pasta"},
							Quantity:   recipe.Quantity{Value: 50, Unit: recipe.QuantityUnitGram},
						},
						{
							Ingredient: recipe.Ingredient{Code: "basil-pesto", Title: "Basil Pesto"},
							Quantity:   recipe.Quantity{Value: 20, Unit: recipe.QuantityUnitGram},
						},
					},
					Instructions: []string{
						"Boil the pasta",
						"Drain the pasta",
						"Combine pesto with pasta",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			ingredientsfile, err := os.Open(tc.IngredientsInput)
			if err != nil {
				t.Fatalf("can't open ingredients file '%s': %v", tc.IngredientsInput, err)
			}
			defer ingredientsfile.Close()

			b, err := recipe.NewBuilderFromYAMLIngredients(ingredientsfile)
			if err != nil {
				t.Fatalf("can't load ingredients from file '%s': %v", tc.IngredientsInput, err)
			}

			for _, recipeInput := range tc.RecipeInputs {
				recipefile, err := os.Open(recipeInput.Path)
				if err != nil {
					t.Fatalf("can't open recipe file '%s': %v", recipeInput, err)
				}
				defer recipefile.Close()

				if err := b.ParseNewYAMLRecipe(recipeInput.Code, recipefile); err != nil {
					t.Errorf("can't load recipe file '%s': %v", recipeInput, err)
				}
			}

			if b.Length() != len(tc.RecipeInputs) {
				t.Fatalf("can't load all recipe. want: %d; got: %d", len(tc.RecipeInputs), b.Length())
			}

			for i := range tc.ExpectedOutput {
				if !reflect.DeepEqual(tc.ExpectedOutput[i], b.Recipes[i]) {
					t.Errorf("invalid recipe parsing.\nwant:\n%#+v\ngot:\n%#+v", tc.ExpectedOutput[i], b.Recipes[i])
				}
			}
		})
	}
}
