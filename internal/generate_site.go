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
	OutputFolderPath string
	PublicURL        string
}

// GenerateSite generates the whole website
func GenerateSite(cfg SiteConfig, ingredients []recipe.Ingredient, recipes []recipe.Recipe) error {
	sitevalues := html.PageSiteValues{
		PublicURL:          cfg.PublicURL,
		CurrentPageSection: html.PageSectionIngredients,
	}

	if err := os.RemoveAll(cfg.OutputFolderPath); err != nil {
		return fmt.Errorf("can't remove output folder '%s': %v", cfg.OutputFolderPath, err)
	}

	if err := copyFolderContent("assets", path.Join(cfg.OutputFolderPath, "assets")); err != nil {
		return fmt.Errorf("can't copy assets folder: %v", err)
	}

	if err := generateIngredients(cfg, sitevalues, ingredients, recipes); err != nil {
		return fmt.Errorf("can't generate ingredients page: %v", err)
	}

	if err := generateRecipes(cfg, sitevalues, recipes); err != nil {
		return fmt.Errorf("can't generate recipes page: %v", err)
	}

	return nil
}

func generateRecipes(cfg SiteConfig, sitevalues html.PageSiteValues, recipes []recipe.Recipe) error {
	type recipesdata struct {
		recipe.Recipe
		TotalDuration   time.Duration
		PricingScale    int
		DifficultyScale int
	}

	data := make([]*recipesdata, len(recipes))
	for i := range recipes {
		data[i] = &recipesdata{
			Recipe:          recipes[i],
			TotalDuration:   time.Duration(recipes[i].Preparation + recipes[i].Cooking),
			PricingScale:    int(recipes[i].Pricing) + 1,
			DifficultyScale: int(recipes[i].Difficulty) + 1,
		}

		err := html.Generate(
			path.Join(cfg.OutputFolderPath, "recipes", recipes[i].Code),
			html.PageRecipeShow,
			html.PageValues{
				Site: sitevalues,
				Data: data[i],
			},
		)
		if err != nil {
			return fmt.Errorf("can't generare recipe '%s': %v", recipes[i].Code, err)
		}
	}

	return html.Generate(
		path.Join(cfg.OutputFolderPath),
		html.PageRecipesList,
		html.PageValues{
			Site: sitevalues,
			Data: data,
		},
	)
}

func generateIngredients(cfg SiteConfig, sitevalues html.PageSiteValues, ingredients []recipe.Ingredient, recipes []recipe.Recipe) error {
	type ingredientdata struct {
		recipe.Ingredient
		NumberOfRecipes int
	}

	data := make(map[string]*ingredientdata, len(ingredients))
	for _, ingredient := range ingredients {
		data[ingredient.Code] = &ingredientdata{Ingredient: ingredient}
	}

	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			data[ingredient.Code].NumberOfRecipes++
		}
	}
	return html.Generate(
		path.Join(cfg.OutputFolderPath, "ingredients"),
		html.PageIngredientsList,
		html.PageValues{
			Site: sitevalues,
			Data: data,
		},
	)
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
