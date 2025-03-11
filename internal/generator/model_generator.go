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
		{{- if isArray .Type }}
	{{ fieldName .Name }} []{{ trimType .Relation }} {{ backtick }}gorm:"-"{{ backtick }}
		{{- else }}
	{{ fieldName (printf "%sID" .Name) }} string {{ backtick }}gorm:"{{ gormTagFK . }}"{{ backtick }}
	{{ fieldName .Name }} *{{ trimType .Relation }} {{ backtick }}gorm:"-"{{ backtick }}
		{{- end }}
	{{- else }}
	{{ fieldName .Name }} {{ goType .Type .IsNonNull }} {{ backtick }}gorm:"{{ gormTag . }}"{{ backtick }}
	{{- end }}
{{- end }}
	CreatedAt time.Time {{ backtick }}gorm:"not null"{{ backtick }}
}
{{ end }}
`
	funcMap := template.FuncMap{
		"fieldName": func(name string) string {
			return strings.Title(name)
		},
		"isArray": func(fieldType string) bool {
			return strings.HasPrefix(fieldType, "[")
		},
		"goType": func(gqlType string, isNonNull bool) string {
			baseType := strings.TrimSuffix(strings.TrimPrefix(gqlType, "["), "]")
			baseType = strings.TrimSuffix(baseType, "!")

			var goType string
			switch baseType {
			case "ID":
				goType = "string"
			case "String":
				goType = "string"
			case "Boolean":
				goType = "bool"
			case "Int":
				goType = "int"
			case "BigInt":
				goType = "string" // Using string for BigInt as it might exceed int64
			case "Date":
				goType = "time.Time"
			default:
				goType = "string"
			}

			if !isNonNull {
				goType = "*" + goType
			}
			return goType
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
			if field.IsUnique && field.IsIndexed {
				tags = append(tags, "uniqueIndex")
			} else if field.IsIndexed {
				tags = append(tags, "index")
			} else if field.IsUnique {
				tags = append(tags, "unique")
			}
			if field.IsNonNull {
				tags = append(tags, "not null")
			}
			return strings.Join(tags, ";")
		},
		"backtick": func() string {
			return "`"
		},
		"trimType": func(t string) string {
			t = strings.TrimPrefix(t, "[")
			t = strings.TrimSuffix(t, "]")
			t = strings.TrimSuffix(t, "!")
			return t
		},
	}
	tpl, err := template.New("models").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}
	return tpl.Execute(f, entities)
}
