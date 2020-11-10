package internal

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

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

	if err := os.RemoveAll(cfg.OutputFolderPath); err != nil {
		return fmt.Errorf("can't remove output folder '%s': %v", cfg.OutputFolderPath, err)
	}

	if err := copyFolderContent("assets", path.Join(cfg.OutputFolderPath, "assets")); err != nil {
		return fmt.Errorf("can't copy assets folder: %v", err)
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
