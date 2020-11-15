package html

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
	"time"
)

// PageSection represents one of the main site area
type PageSection int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=PageSection -linecomment
const (
	// PageSectionIngredients represents the ingredients site area
	PageSectionIngredients PageSection = iota // ingredients
	// PageSectionRecipes represents the recipes site area
	PageSectionRecipes // recipes
)

// Page represents an existing template
type Page int

const (
	// PageIngredientsList represents a page containing the list of ingredients
	PageIngredientsList Page = iota
	// PageRecipesList represents a page containing the list of recipes
	PageRecipesList
	// PageRecipeShow represents a page containing the list of recipes
	PageRecipeShow
)

// PageSiteValues represents the site values of a page
type PageSiteValues struct {
	PublicURL          string
	CurrentPageSection PageSection
}

// PageValues represents a page
type PageValues struct {
	Site  PageSiteValues
	Title string
	Data  interface{}
}

// Renderer represents the HTML page renderer
type Renderer struct {
	templates map[Page]*template.Template
}

// NewRenderer initializes a renderer using a specific template path
func NewRenderer(root string) (*Renderer, error) {
	var r Renderer

	tplNames := map[Page]string{
		PageIngredientsList: "ingredients.html.tmpl",
		PageRecipesList:     "recipes.html.tmpl",
		PageRecipeShow:      "recipe.html.tmpl",
	}

	r.templates = make(map[Page]*template.Template)
	for page, tplName := range tplNames {
		tpl, err := template.
			New(tplName).
			Funcs(helpers()).
			ParseFiles(path.Join(root, "layout.html.tmpl"), path.Join(root, tplName))

		if err != nil {
			return nil, fmt.Errorf("can't load template '%s' from '%s': %v", tplName, root, err)
		}

		r.templates[page] = tpl
	}

	return &r, nil
}

// Generate generates a page from a template
func (r *Renderer) Generate(folder string, page Page, values PageValues) error {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("can't create folder '%s': %v", folder, err)
	}

	f, err := os.OpenFile(path.Join(folder, "index.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create index.html file in '%s': %v", folder, err)
	}
	defer f.Close()

	tpl, ok := r.templates[page]
	if !ok {
		return fmt.Errorf("template doesn't exist")
	}

	if err := tpl.Execute(f, values); err != nil {
		return fmt.Errorf("can't render page: %v", err)
	}

	return nil
}

func helpers() template.FuncMap {
	return template.FuncMap{
		"scale":         displayScale,
		"duration":      displayDuration,
		"ingredientURL": ingredientURL,
		"recipeURL":     recipeURL,
	}
}

func displayDuration(t time.Duration) string {
	return strings.ReplaceAll(strings.ReplaceAll(t.String(), "m0s", "m"), "h0m", "h")
}

func displayScale(icon string, max int, current int) template.HTML {
	var s strings.Builder

	for i := 1; i <= max; i++ {
		cssClass := "uk-text-primary"
		if i > current {
			cssClass = "uk-text-muted"
		}

		fmt.Fprintf(&s, "<span class=\"%s\" uk-icon=\"%s\"></span>", cssClass, icon)
	}

	return template.HTML(s.String())
}

func ingredientURL(publicURL string, code string) string {
	return fmt.Sprintf("%s/ingredients/%s", publicURL, code)
}

func recipeURL(publicURL string, code string) string {
	return fmt.Sprintf("%s/recipes/%s", publicURL, code)
}
