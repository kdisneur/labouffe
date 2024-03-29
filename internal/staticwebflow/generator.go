package staticwebflow

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kdisneur/labouffe/internal/foodaccess"
)

// SiteConfig  represents the site configuration
type SiteConfig struct {
	TemlatesFolderPath string
	OutputFolderPath   string
	PublicURL          string
}

// RecipesView is the data necessary to display a list of recipe template
type RecipesView struct {
	AllDurationRanges []foodaccess.DurationRange
	AllCategories     []foodaccess.Category
	AllDifficulties   []foodaccess.Difficulty
	AllPrices         []foodaccess.Price
	Recipes           []*RecipeView
}

// RecipeView is the data necessary to build a recipe template
type RecipeView struct {
	*foodaccess.Recipe
	Instructions    []*RecipeViewInstruction
	WarningSafeHTML *template.HTML
	TotalDuration   foodaccess.Duration
	PricingScale    int
	DifficultyScale int
}

// RecipeViewInstruction is the instruction data necessary to buld a recipe template
type RecipeViewInstruction struct {
	*foodaccess.Instruction
	WarningSafeHTML *template.HTML
}

// IngredientView is the data necessary to build an ingredient template
type IngredientView struct {
	*foodaccess.Ingredient
	Recipes []*RecipeView
}

// GenerateSite generates the whole website
func GenerateSite(cfg SiteConfig, ingredients []*foodaccess.Ingredient, recipes []*foodaccess.Recipe) error {
	if err := os.RemoveAll(cfg.OutputFolderPath); err != nil {
		return fmt.Errorf("can't remove output folder '%s': %v", cfg.OutputFolderPath, err)
	}

	if err := copyFolderContent(path.Join(cfg.TemlatesFolderPath, "assets"), path.Join(cfg.OutputFolderPath, "assets")); err != nil {
		return fmt.Errorf("can't copy assets folder: %v", err)
	}

	renderer, err := NewRenderer(cfg.TemlatesFolderPath)
	if err != nil {
		return fmt.Errorf("can't initialize html renderer: %v", err)
	}

	durationRanges := foodaccess.AllDurationRanges()
	categories := foodaccess.AllCategories()
	prices := foodaccess.AllPrices()
	difficulties := foodaccess.AllDifficulties()

	recipeviews := make([]*RecipeView, len(recipes))
	for i := range recipes {
		instructions := make([]*RecipeViewInstruction, len(recipes[i].Instructions))
		for j := range recipes[i].Instructions {
			instructions[j] = &RecipeViewInstruction{
				Instruction:     recipes[i].Instructions[j],
				WarningSafeHTML: nl2br(recipes[i].Instructions[j].Warning),
			}
		}

		recipeviews[i] = &RecipeView{
			Recipe:          recipes[i],
			Instructions:    instructions,
			WarningSafeHTML: nl2br(recipes[i].Warning),
			TotalDuration:   foodaccess.Duration(recipes[i].Preparation + +recipes[i].Resting + recipes[i].Cooking),
			PricingScale:    int(recipes[i].Pricing) + 1,
			DifficultyScale: int(recipes[i].Difficulty) + 1,
		}

		err := renderer.Generate(
			path.Join(cfg.OutputFolderPath, "recipes", recipes[i].Code),
			PageRecipeShow,
			PageValues{
				Site: PageSiteValues{
					PublicURL:          cfg.PublicURL,
					CurrentPageSection: PageSectionRecipes,
				},
				Data: recipeviews[i],
			},
		)
		if err != nil {
			return fmt.Errorf("can't generate recipe page '%s': %v", recipes[i].Code, err)
		}
	}

	sort.Slice(recipeviews, func(i int, j int) bool {
		return recipeviews[i].Title < recipeviews[j].Title
	})

	err = renderer.Generate(
		path.Join(cfg.OutputFolderPath),
		PageRecipesList,
		PageValues{
			Site: PageSiteValues{
				PublicURL:          cfg.PublicURL,
				CurrentPageSection: PageSectionRecipes,
			},
			Title: "Les recettes",
			Data: RecipesView{
				Recipes:           recipeviews,
				AllCategories:     categories,
				AllDurationRanges: durationRanges,
				AllDifficulties:   difficulties,
				AllPrices:         prices,
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
			PageRecipesList,
			PageValues{
				Site: PageSiteValues{
					PublicURL:          cfg.PublicURL,
					CurrentPageSection: PageSectionRecipes,
				},
				Title: fmt.Sprintf("%s: Les recettes", ingredient.Title),
				Data: RecipesView{
					Recipes:           ingredient.Recipes,
					AllCategories:     categories,
					AllDurationRanges: durationRanges,
					AllDifficulties:   difficulties,
					AllPrices:         prices,
				},
			},
		)

		if err != nil {
			return fmt.Errorf("can't generate ingredient page '%s': %v", ingredient.Code, err)
		}
	}

	err = renderer.Generate(
		path.Join(cfg.OutputFolderPath, "ingredients"),
		PageIngredientsList,
		PageValues{
			Site: PageSiteValues{
				PublicURL:          cfg.PublicURL,
				CurrentPageSection: PageSectionIngredients,
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

func nl2br(text string) *template.HTML {
	if text == "" {
		return nil
	}

	html := template.HTML(strings.ReplaceAll(template.HTMLEscapeString(text), "\n", "<br>"))

	return &html
}
