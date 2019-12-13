package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var typeNames = []string{
	"decimal128",
	"double",
	"int32",
	"int64",
}

func main() {
	tmpl, err := template.New("mql_template.g4").Funcs(makeFuncMap()).ParseFiles("mql_template.g4")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create("mql.g4")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer out.Close()

	if err = tmpl.Execute(out, nil); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"named_field_any": namedFieldAny,
	}

	for _, typeName := range typeNames {
		typeName := typeName
		funcMap[fmt.Sprintf("named_field_%s", typeName)] = func(name string) string {
			return namedField(name, typeName)
		}
	}

	return funcMap
}

func namedFieldAny(name string) string {
	result := "("
	for i, typeName := range typeNames {
		if i != 0 {
			result += " | "
		}
		result += namedField(name, typeName)
	}
	result += ")"

	return result
}

func namedField(name string, typeName string) string {
	return fmt.Sprintf("TYPE_%s %s %s", strings.ToUpper(typeName), convertNameToMarkers(name), typeName)
}

func convertNameToMarkers(name string) string {
	parts := make([]string, len(name)+1)
	parts[len(name)] = "NUL_BYTE"

	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '$':
			parts[i] = "DOLLAR"
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			parts[i] = strings.ToUpper(string(name[i]))
		default:
			panic(fmt.Sprintf("unknown name character: %s", string(name[i])))
		}
	}

	return strings.Join(parts, " ")
}
