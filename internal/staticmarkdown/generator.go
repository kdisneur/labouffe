package staticmarkdown

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/kdisneur/labouffe/internal/foodaccess"
)

type Config struct {
	TemplatePath     string
	OutputFolderPath string
}

func GenerateMarkdown(cfg Config, ingredients []*foodaccess.Ingredient, recipes []*foodaccess.Recipe) error {
	ingredientsPath := path.Join(cfg.OutputFolderPath, "Ingrédients")
	recipesPath := path.Join(cfg.OutputFolderPath, "Recettes")

	for _, folder := range []string{ingredientsPath, recipesPath} {
		if err := os.RemoveAll(folder); err != nil {
			return fmt.Errorf("can't remove folder %q: %v", folder, err)
		}
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("can't create folder %q: %v", folder, err)
		}
	}

	tplContent, err := os.ReadFile(cfg.TemplatePath)
	if err != nil {
		return fmt.Errorf("can't read template %q: %v", cfg.TemplatePath, err)
	}

	tpl, err := template.New("template").Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("can't parse template %q: %v", cfg.TemplatePath, err)
	}

	for _, ingredient := range ingredients {
		filePath := path.Join(ingredientsPath, fmt.Sprintf("%s.md", ingredient.Title))
		err := os.WriteFile(filePath, []byte("\n"), 0644)
		if err != nil {
			return fmt.Errorf("can't create ingredient file %q: %v", filePath, err)
		}
	}

	for _, recipe := range recipes {
		filePath := path.Join(recipesPath, fmt.Sprintf("%s.md", recipe.Title))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("can't open recipe file %q: %v", filePath, err)
		}

		err = tpl.ExecuteTemplate(file, "template", recipe)
		if err != nil {
			file.Close()
			return fmt.Errorf("can't generate recipe %q: %v", filePath, err)
		}

		file.Close()
	}

	return nil
}
