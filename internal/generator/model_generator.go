package generator

import (
	"os"
	"strings"
	"text/template"
)

func GenerateGoModels(entities []Entity, outputDir string) error {
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
	header := `package models

import (
	"time"
)
`
	if _, err := f.WriteString(header); err != nil {
		return err
	}
	tmpl := `
{{- range .}}
// {{ .Name }} represents the {{ .Name }} entity.
type {{ .Name }} struct {
{{- range .Fields }}
	{{- if .Relation }}
	{{ fieldName (printf "%sID" .Name) }} int {{ backtick }}gorm:"{{ gormTagFK . }}"{{ backtick }}
	{{ fieldName .Name }} *{{ .Relation }} {{ backtick }}gorm:"-"{{ backtick }}
	{{- else }}
	{{ fieldName .Name }} {{ goType .Type .IsNonNull }} {{ backtick }}gorm:"{{ gormTag . }}"{{ backtick }}
	{{- end }}
{{- end }}
}
{{ end }}
`
	funcMap := template.FuncMap{
		"fieldName": func(name string) string {
			return strings.Title(name)
		},
		"goType": func(gqlType string, isNonNull bool) string {
			switch gqlType {
			case "ID":
				return "int"
			case "String":
				return "string"
			case "Boolean":
				return "bool"
			case "Date":
				return "time.Time"
			default:
				return "string"
			}
		},
		"gormTag": func(field Field) string {
			var tags []string
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
		"gormTagFK": func(field Field) string {
			var tags []string
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
		"backtick": func() string {
			return "`"
		},
	}
	tpl, err := template.New("models").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}
	return tpl.Execute(f, entities)
}
