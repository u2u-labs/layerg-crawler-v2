package generator

import (
	"os"
	"strings"
	"text/template"
)

// GenerateGoModels generates Go struct definitions from the parsed entities.
func GenerateGoModels(entities []Entity, outputDir string) error {
	// Create a subdirectory for models so that this file is generated in a separate package.
	modelsDir := outputDir + "/models"
	if err := os.MkdirAll(modelsDir, os.ModePerm); err != nil {
		return err
	}
	filePath := modelsDir + "/models.go"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write package declaration and imports.
	header := `package models

import (
	"time"
)
`
	_, err = f.WriteString(header)
	if err != nil {
		return err
	}

	// Template for each struct.
	tmpl := `
{{- range . }}
// {{ .Name }} represents the {{ .Name }} entity.
type {{ .Name }} struct {
	{{- range .Fields }}
	{{ fieldName .Name }} {{ goType .Type .IsNonNull }} ` + "`gorm:\"{{ gormTag . }}\"`" + `
	{{- end }}
}
{{ end }}
`
	// Template functions for proper naming, type mapping, and tag generation.
	funcMap := template.FuncMap{
		"fieldName": func(name string) string {
			// Capitalize first letter.
			return strings.Title(name)
		},
		"goType": func(gqlType string, isNonNull bool) string {
			// Simple mapping from GraphQL type to Go type.
			switch gqlType {
			case "ID":
				return "string"
			case "String":
				return "string"
			case "Boolean":
				return "bool"
			case "Date":
				return "time.Time"
			default:
				// For relations, assume pointer (adjust as needed).
				return "*" + gqlType
			}
		},
		"gormTag": func(field Field) string {
			// Generate gorm tag.
			tags := []string{}
			if strings.ToLower(field.Name) == "id" {
				tags = append(tags, "primaryKey")
			}
			if field.IsUnique {
				tags = append(tags, "uniqueIndex")
			} else if field.IsIndexed {
				tags = append(tags, "index")
			}
			if field.IsNonNull {
				tags = append(tags, "not null")
			}
			return strings.Join(tags, ";")
		},
	}
	tpl, err := template.New("models").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}
	return tpl.Execute(f, entities)
}
