package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kdisneur/lacuisine/internal/recipe"
)

// Flags represents the set of config flags available for the command line
type Flags struct {
	RecipesFolderPath string
	IngredientsPath   string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	var fcfg Flags

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&fcfg.RecipesFolderPath, "recipes", "./data/recipes", "path to the folder containing all the recipes")
	fs.StringVar(&fcfg.IngredientsPath, "ingredients", "./data/ingredients.yaml", "path to the file containing all the ingredients")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	ingredientfile, err := os.Open(fcfg.IngredientsPath)
	if err != nil {
		return fmt.Errorf("can't open ingredient file '%s': %v", fcfg.IngredientsPath, err)
	}
	defer ingredientfile.Close()

	builder, err := recipe.NewBuilderFromYAMLIngredients(ingredientfile)
	if err != nil {
		return fmt.Errorf("can't import ingredients file '%s': %v", fcfg.IngredientsPath, err)
	}

	err = filepath.Walk(fcfg.RecipesFolderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
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
		return fmt.Errorf("can't import recipes: %v", err)
	}

	for _, recipe := range builder.Recipes {
		fmt.Printf("%#+v\n", recipe)
	}

	return nil
}
