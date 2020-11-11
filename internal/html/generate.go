package html

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
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
type Page struct {
	t *template.Template
}

var (
	// PageIngredientsList represents a page containing the list of ingredients
	PageIngredientsList = &Page{
		t: template.Must(template.New("ingredients.html.tmpl").Funcs(helpers()).ParseFiles("templates/layout.html.tmpl", "templates/ingredients.html.tmpl")),
	}

	// PageRecipesList represents a page containing the list of recipes
	PageRecipesList = &Page{
		t: template.Must(template.New("recipes.html.tmpl").Funcs(helpers()).ParseFiles("templates/layout.html.tmpl", "templates/recipes.html.tmpl")),
	}

	// PageRecipeShow represents a page containing the list of recipes
	PageRecipeShow = &Page{
		t: template.Must(template.New("recipe.html.tmpl").Funcs(helpers()).ParseFiles("templates/layout.html.tmpl", "templates/recipe.html.tmpl")),
	}
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

// Generate generates a page from a template
func Generate(folder string, page *Page, values PageValues) error {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("can't create folder '%s': %v", folder, err)
	}

	f, err := os.OpenFile(path.Join(folder, "index.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create index.html file in '%s': %v", folder, err)
	}
	defer f.Close()

	if err := page.t.Execute(f, values); err != nil {
		return fmt.Errorf("can't render page: %v", err)
	}

	return nil
}

func helpers() template.FuncMap {
	return template.FuncMap{
		"scale": displayScale,
	}
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
