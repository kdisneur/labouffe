package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "more than 2 arguments are required: name and files")
		return
	}

	constantname := os.Args[1]
	filenames := os.Args[2:]

	out, err := os.Create(constantname + ".generated.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create generated file:", err)
		return
	}

	fmt.Fprintln(out, "// Code generated .* DO NOT EDIT.")
	fmt.Fprintln(out, "package", os.Getenv("GOPACKAGE"))
	fmt.Fprintln(out)
	fmt.Fprintf(out, "import (\n")
	fmt.Fprintf(out, "\t\"html/template\"\n")
	fmt.Fprintf(out, "\t\"io\"\n")
	fmt.Fprintf(out, ")\n")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "type %s_Type int\n", constantname)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "const %s %s_Type = 0\n", constantname, constantname)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "var %s_Data = template.New(\"%s\")\n", constantname, constantname)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "func init() {")
	for _, filename := range filenames {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't read source file:", err)
			return
		}

		fmt.Fprintf(out, "\t%s_Data = template.Must(%s_Data.Parse(`%s`))\n", constantname, constantname, content)
	}
	fmt.Fprintln(out, "}")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "func (%s_Type) Render(w io.Writer, data interface{}) error {\n", constantname)
	fmt.Fprintf(out, "\treturn %s_Data.ExecuteTemplate(w, \"layout\", data)\n", constantname)
	fmt.Fprintf(out, "}\n")
}
