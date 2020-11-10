package html

import (
	"fmt"
	"html/template"
	"os"
	"path"
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
type Page string

var tpls = template.Must(template.ParseGlob("templates/*.html.tmpl"))

const (
	// PageIngredientsList represents a page containing the list of ingredients
	PageIngredientsList Page = "ingredients.html.tmpl"
)

// PageSiteValues represents the site values of a page
type PageSiteValues struct {
	PublicURL          string
	CurrentPageSection PageSection
}

// PageValues represents a page
type PageValues struct {
	Site PageSiteValues
	Data interface{}
}

// Generate generates a page from a template
func Generate(folder string, page Page, values PageValues) error {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("can't create folder '%s': %v", folder, err)
	}

	f, err := os.OpenFile(path.Join(folder, "index.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create index.html file in '%s': %v", folder, err)
	}
	defer f.Close()

	if err := tpls.ExecuteTemplate(f, string(page), values); err != nil {
		return fmt.Errorf("can't render page: %v", err)
	}

	return nil
}
