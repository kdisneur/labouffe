package internal_test

import (
	"flag"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/kdisneur/labouffe/internal"
)

var update = flag.Bool("update", false, "update golden files")

func TestGenerateValidSite(t *testing.T) {
	ingredients, recipes, err := internal.LoadIngredientAndRecipes("testdata/valid_recipes/ingredients.yaml", "testdata/valid_recipes/recipes")
	if err != nil {
		t.Fatalf("expected successful ingredients and recipes load: %v", err)
	}

	templatesFolder := "../templates"

	testSiteFolder, err := ioutil.TempDir("", "site")
	if err != nil {
		t.Fatalf("can't generate test folder: %v", err)
	}

	goldenSiteFolder := "testdata/golden_site"
	publicURL := "https://labouffe.local/somefolder"

	if *update {
		err := internal.GenerateSite(internal.SiteConfig{
			TemlatesFolderPath: templatesFolder,
			OutputFolderPath:   goldenSiteFolder,
			PublicURL:          publicURL,
		}, ingredients, recipes)

		if err != nil {
			t.Fatalf("can't update golden files: %v", err)
		}
	}

	err = internal.GenerateSite(internal.SiteConfig{
		TemlatesFolderPath: templatesFolder,
		OutputFolderPath:   testSiteFolder,
		PublicURL:          publicURL,
	}, ingredients, recipes)
	if err != nil {
		t.Fatalf("can't generate site: %v", err)
	}

	var errorDetail strings.Builder
	cmd := exec.Command("diff", "-r", goldenSiteFolder, testSiteFolder)
	cmd.Stdout = &errorDetail
	cmd.Stderr = &errorDetail
	if err := cmd.Run(); err != nil {
		t.Fatalf("site folder are different: %v\nDetails:\n%s", err, errorDetail.String())
	}
}
