package generator

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

// GenerateMigrationScripts generates a migration script from the schema.
func GenerateMigrationScripts(entities []Entity, outputDir string) error {
	// Create a migration file name based on current timestamp.
	timestamp := time.Now().Format("20060102150405")
	filePath := fmt.Sprintf("%s/migration_%s.sql", outputDir, timestamp)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := `-- +goose Up
-- Migration script generated from GraphQL schema

`
	_, err = f.WriteString(header)
	if err != nil {
		return err
	}

	// Template to generate CREATE TABLE statements.
	tmpl := `
{{- range . }}
{{$fields := .Fields}}
CREATE TABLE {{ .Name | toSnakeCase }} (
	{{- range $index, $field := $fields }}
	{{ $field.Name | toSnakeCase }} {{ sqlType $field.Type }} {{ if $field.IsNonNull }}NOT NULL{{ end }}{{ if not (isLast $index $fields) }},{{ end }}
	{{- end }}
);
{{ end }}
`
	funcMap := template.FuncMap{
		"toSnakeCase": func(name string) string {
			// Naïve conversion to lower-case (improve as needed).
			return strings.ToLower(name)
		},
		"sqlType": func(gqlType string) string {
			// Map GraphQL type to SQL types.
			switch gqlType {
			case "ID", "String":
				return "VARCHAR"
			case "Boolean":
				return "BOOLEAN"
			case "Date":
				return "TIMESTAMP"
			default:
				// For relation types, use VARCHAR (or adjust with proper references).
				return "VARCHAR"
			}
		},
		"isLast": func(index int, fields []Field) bool {
			return index == len(fields)-1
		},
	}
	tpl, err := template.New("migration").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}
	if err := tpl.Execute(f, entities); err != nil {
		return err
	}

	// Footer for migration down (this is a stub; adjust as needed).
	downFooter := `
-- +goose Down
`
	_, err = f.WriteString(downFooter)
	return err
}
