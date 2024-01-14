package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/kdisneur/labouffe/internal/foodaccess"
	"github.com/kdisneur/labouffe/internal/staticmarkdown"
	"github.com/kdisneur/labouffe/internal/staticwebflow"
)

// Flags represents the set of config flags available for the command line
type Flags struct {
	RecipesFolderPath string
	IngredientsPath   string
	OutputFolderPath  string
	DeveloperMode     bool
	DeveloperModePort int
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
	fs.BoolVar(&fcfg.DeveloperMode, "dev", false, "start a server to render HTML pages")
	fs.IntVar(&fcfg.DeveloperModePort, "dev-http-port", 8080, "http port of the developer mode server")
	fs.StringVar(&fcfg.RecipesFolderPath, "recipes", "./data/recipes", "path to the folder containing all the recipes")
	fs.StringVar(&fcfg.IngredientsPath, "ingredients", "./data/ingredients.yaml", "path to the file containing all the ingredients")
	fs.StringVar(&fcfg.OutputFolderPath, "output", "./public", "path to the generated site")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	ingredients, recipes, err := foodaccess.LoadIngredientAndRecipes(fcfg.IngredientsPath, fcfg.RecipesFolderPath)
	if err != nil {
		return fmt.Errorf("can't load data: %v", err)
	}

	if err := os.MkdirAll(fcfg.OutputFolderPath, 0755); err != nil {
		return fmt.Errorf("can't create output folder '%s': %v", fcfg.OutputFolderPath, err)
	}

	if fs.Arg(0) == "export" {
		cfg := staticmarkdown.Config{
			TemplatePath:     "templates/obsidian.md",
			OutputFolderPath: fcfg.OutputFolderPath,
		}
		if err := staticmarkdown.GenerateMarkdown(cfg, ingredients, recipes); err != nil {
			return fmt.Errorf("can't export markdown: %v", err)
		}
	} else {
		sitecfg := staticwebflow.SiteConfig{
			TemlatesFolderPath: "templates",
			OutputFolderPath:   fcfg.OutputFolderPath,
		}
		if err := staticwebflow.GenerateSite(sitecfg, ingredients, recipes); err != nil {
			return fmt.Errorf("can't generate site: %v", err)
		}

		if fcfg.DeveloperMode {
			fmt.Printf("start developer server: http://127.0.0.1:%d\n", fcfg.DeveloperModePort)
			return http.ListenAndServe(fmt.Sprintf(":%d", fcfg.DeveloperModePort), http.FileServer(http.Dir(fcfg.OutputFolderPath)))
		}
	}

	return nil
}
