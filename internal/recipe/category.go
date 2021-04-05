package recipe

import (
	"fmt"
	"strings"
)

// Category represents a recipe category
type Category int

//go:generate go run ../../vendor/golang.org/x/tools/cmd/stringer/stringer.go -type=Category -linecomment
const (
	// CategoryBiscuit represents a biscuit kind of recipe
	CategoryBiscuit Category = iota // biscuit
	// CategoryCake represents a cake kind of recipe
	CategoryCake // gâteau
	// CategoryColdDishes represents a cold dishes kind of recipe
	CategoryColdDishes // plat froid
	// CategoryHotDishes represents a hot dishes kind of recipe
	CategoryHotDishes // plat chaud
	// CategorySideDishes represents a side dishes kind of recipe
	CategorySideDishes // entrée
	// CategorySauce represents a sauce kind of recipe
	CategorySauce // sauce
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
