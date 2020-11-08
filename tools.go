// +build tools

package tools

import (
	// ensures the code follow the language best practice
	_ "golang.org/x/lint/golint"
	// generates a String method for the custom enums
	_ "golang.org/x/tools/cmd/stringer"
	// performs static analysis on the code base
	_ "honnef.co/go/tools/cmd/staticcheck"
)
