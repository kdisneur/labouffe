package internal

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/kdisneur/labouffe/internal/html"
	"github.com/kdisneur/labouffe/internal/recipe"
)

// SiteConfig  represents the site configuration
type SiteConfig struct {
	TemlatesFolderPath string
	OutputFolderPath   string
	PublicURL          string
}

// RecipesView is the data necessary to display a list of recipe template
type RecipesView struct {
	AllCategories   []recipe.Category
	AllDifficulties []recipe.Difficulty
	AllPrices       []recipe.Price
	Recipes         []*RecipeView
}

// RecipeView is the data necessary to build a recipe template
type RecipeView struct {
	*recipe.Recipe
	TotalDuration   time.Duration
	PricingScale    int
	DifficultyScale int
}

// IngredientView is the data necessary to build an ingredient template
type IngredientView struct {
	*recipe.Ingredient
	Recipes []*RecipeView
}

// GenerateSite generates the whole website
func GenerateSite(cfg SiteConfig, ingredients []*recipe.Ingredient, recipes []*recipe.Recipe) error {
	if err := os.RemoveAll(cfg.OutputFolderPath); err != nil {
		return fmt.Errorf("can't remove output folder '%s': %v", cfg.OutputFolderPath, err)
	}

	if err := copyFolderContent(path.Join(cfg.TemlatesFolderPath, "assets"), path.Join(cfg.OutputFolderPath, "assets")); err != nil {
		return fmt.Errorf("can't copy assets folder: %v", err)
	}

	renderer, err := html.NewRenderer(cfg.TemlatesFolderPath)
	if err != nil {
		return fmt.Errorf("can't initialize html renderer: %v", err)
	}

	categories := recipe.AllCategories()
	prices := recipe.AllPrices()
	difficulties := recipe.AllDifficulties()

	recipeviews := make([]*RecipeView, len(recipes))
	for i := range recipes {
		recipeviews[i] = &RecipeView{
			Recipe:          recipes[i],
			TotalDuration:   time.Duration(recipes[i].Preparation + recipes[i].Cooking),
			PricingScale:    int(recipes[i].Pricing) + 1,
			DifficultyScale: int(recipes[i].Difficulty) + 1,
		}

		err := renderer.Generate(
			path.Join(cfg.OutputFolderPath, "recipes", recipes[i].Code),
			html.PageRecipeShow,
			html.PageValues{
				Site: html.PageSiteValues{
					PublicURL:          cfg.PublicURL,
					CurrentPageSection: html.PageSectionRecipes,
				},
				Data: recipeviews[i],
			},
		)
		if err != nil {
			return fmt.Errorf("can't generate recipe page '%s': %v", recipes[i].Code, err)
		}
	}

	err = renderer.Generate(
		path.Join(cfg.OutputFolderPath),
		html.PageRecipesList,
		html.PageValues{
			Site: html.PageSiteValues{
				PublicURL:          cfg.PublicURL,
				CurrentPageSection: html.PageSectionRecipes,
			},
			Title: "Les recettes",
			Data: RecipesView{
				Recipes:         recipeviews,
				AllCategories:   categories,
				AllDifficulties: difficulties,
				AllPrices:       prices,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("can't generate all recipes page: %v", err)
	}

	data := make(map[string]*IngredientView, len(ingredients))
	for _, ingredient := range ingredients {
		data[ingredient.Code] = &IngredientView{Ingredient: ingredient}
	}

	for _, recipe := range recipeviews {
		for _, ingredient := range recipe.Ingredients {
			data[ingredient.Code].Recipes = append(data[ingredient.Code].Recipes, recipe)
		}
	}

	for _, ingredient := range data {
		err := renderer.Generate(
			path.Join(cfg.OutputFolderPath, "ingredients", ingredient.Code),
			html.PageRecipesList,
			html.PageValues{
				Site: html.PageSiteValues{
					PublicURL:          cfg.PublicURL,
					CurrentPageSection: html.PageSectionRecipes,
				},
				Title: fmt.Sprintf("%s: Les recettes", ingredient.Title),
				Data: RecipesView{
					Recipes:         ingredient.Recipes,
					AllCategories:   categories,
					AllDifficulties: difficulties,
					AllPrices:       prices,
				},
			},
		)

		if err != nil {
			return fmt.Errorf("can't generate ingredient page '%s': %v", ingredient.Code, err)
		}
	}

	err = renderer.Generate(
		path.Join(cfg.OutputFolderPath, "ingredients"),
		html.PageIngredientsList,
		html.PageValues{
			Site: html.PageSiteValues{
				PublicURL:          cfg.PublicURL,
				CurrentPageSection: html.PageSectionIngredients,
			},
			Title: "Les ingrédients",
			Data:  data,
		},
	)
	if err != nil {
		return fmt.Errorf("can't generate all ingredients page: %v", err)
	}

	return nil
}

func copyFolderContent(source string, destination string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("can't compute destination path for '%s': %v", path, err)
		}

		destinationPath := filepath.Join(destination, relpath)

		if info.IsDir() {
			return os.MkdirAll(destinationPath, 0755)
		}

		srcfile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("can't open source file '%s': %v", path, err)
		}

		dstfile, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("can't open destination file '%s': %v", destinationPath, err)
		}

		if _, err := io.Copy(dstfile, srcfile); err != nil {
			return fmt.Errorf("can't copy file from '%s' to '%s': %v", path, destinationPath, err)
		}

		return nil
	})
}
