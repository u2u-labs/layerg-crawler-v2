package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

// GenerateMigrationScripts generates a migration script from the schema.
func GenerateMigrationScripts(entities []Entity, outputDir string) error {
	// Create a dedicated subfolder for generated migrations.
	migrationsDir := outputDir + "/migrations"
	if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
		return err
	}

	// Load previous schema snapshot if it exists.
	snapshotFile := migrationsDir + "/schema_snapshot.json"
	var prevEntities []Entity
	if data, err := os.ReadFile(snapshotFile); err == nil {
		if err := json.Unmarshal(data, &prevEntities); err != nil {
			return err
		}
	}

	var migrationSQL string
	if len(prevEntities) > 0 {
		// Generate an incremental (diff) migration.
		diffSQL, err := generateDiffMigration(prevEntities, entities)
		if err != nil {
			return err
		}
		migrationSQL = diffSQL
	} else {
		// No previous snapshot exists; generate full migration.
		fullSQL, err := generateFullMigration(entities)
		if err != nil {
			return err
		}
		migrationSQL = fullSQL
	}

	// Ensure migrationSQL is not empty so that a migration file is generated.
	if strings.TrimSpace(migrationSQL) == "" {
		migrationSQL = "-- No schema changes detected\n"
	}

	timestamp := time.Now().Format("20060102150405")
	filePath := fmt.Sprintf("%s/%s_migration.sql", migrationsDir, timestamp)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header := `-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

`
	if _, err := f.WriteString(header); err != nil {
		return err
	}
	if _, err := f.WriteString(migrationSQL); err != nil {
		return err
	}
	downFooter := `
-- +goose Down
`
	if _, err := f.WriteString(downFooter); err != nil {
		return err
	}

	// Update snapshot with current schema.
	snapshotData, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(snapshotFile, snapshotData, 0644); err != nil {
		return err
	}
	return nil
}

// generateFullMigration generates a full migration SQL using our template.
func generateFullMigration(entities []Entity) (string, error) {
	tmpl := `
{{- range . }}
{{$fields := .Fields}}
CREATE TABLE "{{ .Name | toSnakeCase }}" (
	{{- range $index, $field := $fields }}
	"{{ $field.Name | toSnakeCase }}" {{ sqlType $field.Type }} {{ if $field.IsNonNull }}NOT NULL{{ end }}{{ if not (isLast $index $fields) }},{{ end }}
	{{- end }}
);
{{ end }}
`
	funcMap := template.FuncMap{
		"toSnakeCase": func(name string) string {
			return strings.ToLower(name)
		},
		"sqlType": func(gqlType string) string {
			switch gqlType {
			case "ID", "String":
				return "VARCHAR"
			case "Boolean":
				return "BOOLEAN"
			case "Date":
				return "TIMESTAMP"
			default:
				return "VARCHAR"
			}
		},
		"isLast": func(index int, fields []Field) bool {
			return index == len(fields)-1
		},
	}
	tpl, err := template.New("fullMigration").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := tpl.Execute(&sb, entities); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// generateDiffMigration generates a migration SQL containing changes between previous and current schema.
// For simplicity, here we only generate CREATE TABLE statements for new entities.
func generateDiffMigration(prev, curr []Entity) (string, error) {
	var sb strings.Builder
	for _, newEntity := range curr {
		exists := false
		for _, oldEntity := range prev {
			if strings.EqualFold(newEntity.Name, oldEntity.Name) {
				exists = true
				break
			}
		}
		if !exists {
			full, err := generateFullMigration([]Entity{newEntity})
			if err != nil {
				return "", err
			}
			sb.WriteString(full)
			sb.WriteString("\n")
		}
	}
	// (Optional) Compare existing entities and generate ALTER TABLE statements for new fields.
	// This example only handles new entities.
	return sb.String(), nil
}
