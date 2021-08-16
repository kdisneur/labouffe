package foodaccess_test

import (
	"testing"
	"time"

	"github.com/kdisneur/labouffe/internal/foodaccess"
)

func TestLoadingValidRecipes(t *testing.T) {
	ingredients, recipes, err := foodaccess.LoadIngredientAndRecipes("testdata/valid_recipes/ingredients.yaml", "testdata/valid_recipes/recipes")
	if err != nil {
		t.Fatalf("expected successful ingredients and recipes load: %v", err)
	}

	expectedIngredients := []*foodaccess.Ingredient{
		{Code: "farine", Title: "Farine"},
		{Code: "oeuf", Title: "Oeuf"},
		{Code: "sucre-glace", Title: "Sucre glace"},
		{Code: "puree", Title: "Purée"},
		{Code: "recipe4-no-override", Title: "Recipe 4"},
	}

	expectedRecipes := []*foodaccess.Recipe{
		{
			Code:        "recipe-1-biscuit-facile-economique",
			Title:       "Recipe 1",
			Category:    foodaccess.CategoryBiscuit,
			Preparation: foodaccess.Duration(10 * time.Minute),
			Cooking:     foodaccess.Duration(20 * time.Minute),
			Difficulty:  foodaccess.DifficultyEasy,
			Pricing:     foodaccess.PriceCheap,
			Servings:    foodaccess.Servings{Quantity: 20, Type: foodaccess.ServingsTypeUnits},
			Ingredients: []*foodaccess.IncludedIngredient{
				{
					Ingredient: expectedIngredients[1],
					Quantity:   foodaccess.Quantity{Value: 1, Unit: foodaccess.QuantityUnitNoUnit},
					Details:    "jaune",
				},
				{
					Ingredient: expectedIngredients[0],
					Quantity:   foodaccess.Quantity{Value: 300, Unit: foodaccess.QuantityUnitGram},
				},
				{
					Ingredient: expectedIngredients[2],
				},
			},
			Instructions: []*foodaccess.Instruction{
				{Title: "La première instruction"},
				{Title: "La seconde instruction"},
				{Title: "La troisième instruction"},
			},
		},
		{
			Code:        "recipe-2-plat-froid-moyen-abordable",
			Title:       "Recipe 2",
			Category:    foodaccess.CategoryColdMeal,
			Preparation: foodaccess.Duration(10 * time.Minute),
			Resting:     foodaccess.Duration(3 * time.Minute),
			Cooking:     foodaccess.Duration(20 * time.Minute),
			Difficulty:  foodaccess.DifficultyAverage,
			Pricing:     foodaccess.PriceAffordable,
			Servings:    foodaccess.Servings{Quantity: 4, Type: foodaccess.ServingsTypeGuests},
			Ingredients: []*foodaccess.IncludedIngredient{
				{
					Ingredient: expectedIngredients[1],
					Quantity:   foodaccess.Quantity{Value: 3, Unit: foodaccess.QuantityUnitNoUnit},
				},
				{
					Ingredient: expectedIngredients[0],
					Quantity:   foodaccess.Quantity{Value: 1.5, Unit: foodaccess.QuantityUnitKilogram},
				},
			},
			Instructions: []*foodaccess.Instruction{
				{Title: "La première instruction"},
				{Title: "La seconde instruction"},
				{Title: "La troisième instruction"},
			},
		},
		{
			Code:        "recipe-3-plat-chaud-difficile-cher",
			Title:       "Recipe 3",
			Category:    foodaccess.CategoryHotMeal,
			Preparation: foodaccess.Duration(45 * time.Minute),
			Cooking:     foodaccess.Duration(time.Hour + 14*time.Minute),
			Difficulty:  foodaccess.DifficultyHard,
			Pricing:     foodaccess.PriceExpensive,
			Servings:    foodaccess.Servings{Quantity: 8, Type: foodaccess.ServingsTypeGuests},
			Ingredients: []*foodaccess.IncludedIngredient{
				{
					Ingredient: expectedIngredients[1],
					Quantity:   foodaccess.Quantity{Value: 3, Unit: foodaccess.QuantityUnitNoUnit},
				},
				{
					Ingredient:   expectedIngredients[0],
					Quantity:     foodaccess.Quantity{Value: 1.5, Unit: foodaccess.QuantityUnitKilogram},
					Alternatives: []*foodaccess.Ingredient{expectedIngredients[2]},
				},
				{
					Ingredient: expectedIngredients[3],
					Quantity:   foodaccess.Quantity{Value: 100, Unit: foodaccess.QuantityUnitGram},
				},
			},
			Instructions: []*foodaccess.Instruction{
				{Title: "La première instruction"},
				{Title: "La seconde instruction"},
				{Title: "La troisième instruction"},
			},
		},
		{
			Code:        "recipe-4-plat-chaud-facile-economique",
			Title:       "Recipe 4",
			Category:    foodaccess.CategoryHotMeal,
			Preparation: foodaccess.Duration(10 * time.Minute),
			Cooking:     foodaccess.Duration(20 * time.Minute),
			Difficulty:  foodaccess.DifficultyEasy,
			Pricing:     foodaccess.PriceCheap,
			Servings:    foodaccess.Servings{Quantity: 2, Type: foodaccess.ServingsTypeGuests},
			Ingredients: []*foodaccess.IncludedIngredient{
				{
					Ingredient: expectedIngredients[1],
					Quantity:   foodaccess.Quantity{Value: 1, Unit: foodaccess.QuantityUnitNoUnit},
				},
				{
					Ingredient: expectedIngredients[0],
					Quantity:   foodaccess.Quantity{Value: 10, Unit: foodaccess.QuantityUnitGram},
				},
			},
			Instructions: []*foodaccess.Instruction{
				{Title: "La première instruction"},
				{Title: "La seconde instruction"},
				{
					Title:   "La troisième instruction",
					Warning: "mais attention de bien la faire avant la seconde",
				},
			},
		},
		{
			Code:        "recipe-5-entree-facile-economique",
			Title:       "Recipe 5",
			Category:    foodaccess.CategoryStarterDish,
			Warning:     "Ceci est une recette a préparer à l'avance",
			Preparation: foodaccess.Duration(10 * time.Minute),
			Cooking:     foodaccess.Duration(35 * time.Minute),
			Difficulty:  foodaccess.DifficultyEasy,
			Pricing:     foodaccess.PriceCheap,
			Servings:    foodaccess.Servings{Quantity: 2, Type: foodaccess.ServingsTypeGuests},
			Ingredients: []*foodaccess.IncludedIngredient{
				{
					Ingredient: expectedIngredients[1],
					Quantity:   foodaccess.Quantity{Value: 1, Unit: foodaccess.QuantityUnitNoUnit},
				},
				{
					Ingredient: expectedIngredients[0],
					Quantity:   foodaccess.Quantity{Value: 12, Unit: foodaccess.QuantityUnitGram},
				},
			},
			Instructions: []*foodaccess.Instruction{
				{Title: "La première instruction"},
				{Title: "La seconde instruction"},
				{Title: "La troisième instruction"},
			},
		},
	}

	expectedIngredients[3].Recipe = expectedRecipes[3]
	expectedIngredients[4].Recipe = expectedRecipes[3]

	expectedRecipes[2].Instructions[0].Recipe = expectedRecipes[3]

	if len(expectedIngredients) != len(ingredients) {
		t.Fatalf("wrong number of ingredients. want: %d; have: %d", len(expectedIngredients), len(ingredients))
	}

	if len(expectedRecipes) != len(recipes) {
		t.Fatalf("wrong number of recipes. want: %d; have: %d", len(expectedRecipes), len(recipes))
	}

	for i := range expectedIngredients {
		assertIngredient(t, expectedIngredients[i], ingredients[i])
	}

	for i := range expectedRecipes {
		assertRecipe(t, expectedRecipes[i], recipes[i])
	}
}

func assertIngredient(t *testing.T, want *foodaccess.Ingredient, got *foodaccess.Ingredient) {
	t.Helper()

	if want != nil && got == nil {
		t.Errorf("expected ingredient '%s' to be set got nil", want.Code)
	}

	if want == nil && got != nil {
		t.Errorf("expected ingredient to be nil; got: %v", got.Code)
	}

	if want == nil && got == nil {
		return
	}

	if want.Code != got.Code {
		t.Errorf("unexpected ingredient code: want: %s; got: %s", want.Code, got.Code)
	}

	if want.Title != got.Title {
		t.Errorf("unexpected ingredient '%s' title: want: %s; got: %s", want.Code, want.Title, got.Title)
	}

	assertRecipe(t, want.Recipe, got.Recipe)
}

func assertRecipe(t *testing.T, want *foodaccess.Recipe, got *foodaccess.Recipe) {
	t.Helper()

	if want != nil && got == nil {
		t.Fatalf("expected recipe '%s' to be set got nil", want.Code)
	}

	if want == nil && got != nil {
		t.Fatalf("expected recipe to be nil; got: %v", got.Code)
	}

	if want == nil && got == nil {
		return
	}

	if want.Code != got.Code {
		t.Errorf("unexpected recipe code: want: %s; got: %s", want.Code, got.Code)
	}

	if want.Title != got.Title {
		t.Errorf("unexpected recipe '%s' title: want: %s; got: %s", want.Code, want.Title, got.Title)
	}

	if want.Warning != got.Warning {
		t.Errorf("unexpected recipe '%s' warning: want: %s; got: %s", want.Code, want.Warning, got.Warning)
	}

	if want.Category != got.Category {
		t.Errorf("unexpected recipe '%s' category: want: %s; got: %s", want.Code, want.Category, got.Category)
	}

	if want.Preparation != got.Preparation {
		t.Errorf("unexpected recipe '%s' preparation: want: %s; got: %s", want.Code, want.Preparation, got.Preparation)
	}

	if want.Resting != got.Resting {
		t.Errorf("unexpected recipe '%s' resting: want: %s; got: %s", want.Code, want.Resting, got.Resting)
	}

	if want.Cooking != got.Cooking {
		t.Errorf("unexpected recipe '%s' cooking: want: %s; got: %s", want.Code, want.Cooking, got.Cooking)
	}

	if want.Difficulty != got.Difficulty {
		t.Errorf("unexpected recipe '%s' difficulty: want: %s; got: %s", want.Code, want.Difficulty, got.Difficulty)
	}

	if want.Pricing != got.Pricing {
		t.Errorf("unexpected recipe '%s' pricing: want: %s; got: %s", want.Code, want.Pricing, got.Pricing)
	}

	if want.Servings != got.Servings {
		t.Errorf("unexpected recipe '%s' servings: want: %#+v; got: %#+v", want.Code, want.Servings, got.Servings)
	}

	if len(want.Ingredients) != len(got.Ingredients) {
		t.Fatalf("unexpected recipe '%s' ingredients number: want: %d; got: %d", want.Code, len(want.Ingredients), len(got.Ingredients))
	}

	for i := range want.Ingredients {
		wantIngredient := want.Ingredients[i]
		gotIngredient := want.Ingredients[i]

		if wantIngredient.Code != gotIngredient.Code {
			t.Errorf("unexpected recipe '%s' ingredient code: want: %s; got: %s", want.Code, wantIngredient.Code, gotIngredient.Code)
		}

		if wantIngredient.Title != gotIngredient.Title {
			t.Errorf("unexpected recipe '%s' ingredient title: want: %s; got: %s", want.Code, wantIngredient.Title, gotIngredient.Title)
		}

		if wantIngredient.Quantity != gotIngredient.Quantity {
			t.Errorf("unexpected recipe '%s' ingredient quantity: want: %#+v; got: %#+vs", want.Code, wantIngredient.Quantity, gotIngredient.Quantity)
		}

		if wantIngredient.Details != gotIngredient.Details {
			t.Errorf("unexpected recipe '%s' ingredient details: want: %s; got: %s", want.Code, wantIngredient.Details, gotIngredient.Details)
		}

		t.Logf("asserting included ingredient '%s' for recipe '%s'", wantIngredient.Code, want.Code)
		assertIngredient(t, wantIngredient.Ingredient, gotIngredient.Ingredient)

		if len(wantIngredient.Alternatives) != len(gotIngredient.Alternatives) {
			t.Fatalf("unexpected recipe '%s' ingredients alternative number for '%s': want: %d; got: %d", want.Code, wantIngredient.Code, len(wantIngredient.Alternatives), len(gotIngredient.Alternatives))
		}

		for j := range wantIngredient.Alternatives {
			t.Logf("asserting included ingredient alternative '%s' for recipe '%s'", wantIngredient.Alternatives[j].Code, want.Code)
			assertIngredient(t, wantIngredient.Alternatives[j], gotIngredient.Alternatives[j])
		}
	}

	if len(want.Instructions) != len(got.Instructions) {
		t.Fatalf("unexpected recipe '%s' instructions number: want: %d; got: %d", want.Code, len(want.Instructions), len(got.Instructions))
	}

	for i := range want.Instructions {
		wantInstruction := want.Instructions[i]
		gotInstruction := got.Instructions[i]

		if wantInstruction.Title != gotInstruction.Title {
			t.Errorf("unexpected recipe '%s' instruction title: want: %s; got: %s", want.Code, wantInstruction.Title, gotInstruction.Title)
		}

		if wantInstruction.Recipe != nil && gotInstruction.Recipe == nil {
			t.Errorf("unexpected recipe '%s' instruction recipe: want a linked recipe but did not get one", want.Code)
		}

		if wantInstruction.Recipe == nil && gotInstruction.Recipe != nil {
			t.Errorf("unexpected recipe '%s' instruction recipe: didn't want a linked recipe but got one", want.Code)
		}

		if wantInstruction.Recipe != nil {
			t.Logf("asserting likned reciped to an instruction for recipe '%s'", want.Code)
			assertRecipe(t, wantInstruction.Recipe, gotInstruction.Recipe)
		}

		if wantInstruction.Warning != gotInstruction.Warning {
			t.Errorf("unexpected recipe '%s' instruction warning: want: %s; got: %s", want.Code, wantInstruction.Warning, gotInstruction.Warning)
		}

		t.Logf("asserting instruction recipe for '%s'", want.Code)
		assertRecipe(t, wantInstruction.Recipe, gotInstruction.Recipe)
	}
}
