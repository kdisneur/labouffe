package foodaccess

import (
	"fmt"
	"strings"
)

// Category represents a recipe category
type Category int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=Category -linecomment
const (
	// CategoryColdMeal represents a cold meal kind of recipe
	CategoryColdMeal Category = iota // plat froid
	// CategoryHotMeal represents a hot meal kind of recipe
	CategoryHotMeal // plat chaud
	// CategoryEggDish represents an egg dish kind of recipe
	CategoryEggDish // oeuf
	// CategoryMeatDish represents a meat dish kind of recipe
	CategoryMeatDish // viande
	// CategoryFishDish represents a fish dish kind of recipe
	CategoryFishDish // poisson
	// CategorySideDishes represents a side dish kind of recipe
	CategorySideDish // accompagnement
	// CategoryBiscuit represents a biscuit kind of recipe
	CategoryBiscuit // biscuit
	// CategoryDessert represents a dessert kind of recipe
	CategoryDessert // dessert
	// CategoryStarterDish represents a starter dish kind of recipe
	CategoryStarterDish // entrée
	// CategorySauce represents a sauce kind of recipe
	CategorySauce // sauce
	// CategorySpice represents a spice kind of recipe
	CategorySpice // épice
)

// AllCategories returns the full list of categories
func AllCategories() []Category {
	var categories []Category
	for i := 0; i < len(_Category_index)-1; i++ {
		categories = append(categories, Category(i))
	}

	return categories
}

// UnmarshalYAML is the function in charge of unmarshalling the string value to a Go constant
func (c *Category) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var namedCategory string
	if err := unmarshal(&namedCategory); err != nil {
		return fmt.Errorf("can't unmarshal named category, expected a string: %v", err)
	}

	var allNamedCategory []string
	for i := 0; i < len(_Category_index)-1; i++ {
		allNamedCategory = append(allNamedCategory, _Category_name[_Category_index[i]:_Category_index[i+1]])
	}

	for i := range allNamedCategory {
		if allNamedCategory[i] == strings.ToLower(namedCategory) {
			*c = Category(i)

			return nil
		}
	}

	return fmt.Errorf("unsupported category '%s'. valid category are: %s", namedCategory, strings.Join(allNamedCategory, ", "))
}
