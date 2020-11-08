package html

import (
	"fmt"
	"os"
	"path"
)

// Generate generates a dummy hello world page
func Generate(folder string) error {
	stat, err := os.Stat(folder)
	if err != nil {
		return fmt.Errorf("can't get stat from output folder '%s': %v", folder, err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("output folder '%s' is not a folder", folder)
	}

	f, err := os.OpenFile(path.Join(folder, "index.html"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can't create index.html file: %v", err)
	}
	defer f.Close()

	f.WriteString(`<html><body>Hello World</body></html>`)

	return nil
}
