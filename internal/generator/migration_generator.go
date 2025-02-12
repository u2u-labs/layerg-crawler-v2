package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func GenerateMigrationScripts(entities []Entity, outputDir string) error {
	migrationsDir := outputDir + "/migrations"
	if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
		return err
	}
	snapshotFile := migrationsDir + "/schema_snapshot.json"
	var prevEntities []Entity
	if data, err := os.ReadFile(snapshotFile); err == nil {
		if err := json.Unmarshal(data, &prevEntities); err != nil {
			return err
		}
	}
	var migrationSQL string
	if len(prevEntities) > 0 {
		diffSQL, err := generateDiffMigration(prevEntities, entities)
		if err != nil {
			return err
		}
		migrationSQL = diffSQL
	} else {
		fullSQL, err := generateFullMigration(entities)
		if err != nil {
			return err
		}
		migrationSQL = fullSQL
	}
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
	snapshotData, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(snapshotFile, snapshotData, 0644); err != nil {
		return err
	}
	return nil
}

func generateFullMigration(entities []Entity) (string, error) {
	sortedEntities, err := sortEntities(entities)
	if err != nil {
		return "", err
	}
	// Build a map from entity name to its primary key type.
	pkMap := make(map[string]string)
	for _, entity := range sortedEntities {
		for _, field := range entity.Fields {
			if strings.ToLower(field.Name) == "id" && field.Relation == "" {
				switch strings.ToLower(field.Type) {
				case "id":
					pkMap[entity.Name] = "SERIAL"
				case "bigint":
					pkMap[entity.Name] = "DECIMAL"
				default:
					pkMap[entity.Name] = "SERIAL"
				}
				break
			}
		}
	}
	var sb strings.Builder
	for _, entity := range sortedEntities {
		tableName := toSnakeCase(entity.Name)
		sb.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))
		var colDefs []string
		var foreignKeys []string
		for _, field := range entity.Fields {
			var colName string
			if field.Relation != "" {
				colName = toSnakeCase(field.Name) + "_id"
			} else {
				colName = toSnakeCase(field.Name)
			}
			var colType string
			if field.Relation != "" {
				// Lookup the referenced entity's primary key type.
				if pk, ok := pkMap[field.Relation]; ok {
					if pk == "DECIMAL" {
						colType = "DECIMAL"
					} else {
						colType = "INTEGER"
					}
				} else {
					colType = "INTEGER"
				}
			} else {
				switch strings.ToLower(field.Type) {
				case "id":
					colType = "SERIAL"
				case "bigint":
					colType = "DECIMAL"
				case "string":
					colType = "TEXT"
				case "boolean":
					colType = "BOOLEAN"
				case "date":
					colType = "TIMESTAMPTZ"
				default:
					colType = "TEXT"
				}
			}
			colDef := fmt.Sprintf("    \"%s\" %s", colName, colType)
			if field.IsNonNull && strings.ToLower(field.Type) != "id" {
				colDef += " NOT NULL"
			}
			if strings.ToLower(field.Type) == "id" && field.Relation == "" {
				colDef += " PRIMARY KEY"
			}
			colDefs = append(colDefs, colDef)
			if field.Relation != "" {
				refTable := toSnakeCase(field.Relation)
				foreignKeys = append(foreignKeys, fmt.Sprintf("    FOREIGN KEY (\"%s\") REFERENCES \"%s\"(\"id\")", colName, refTable))
			}
		}
		allDefs := strings.Join(colDefs, ",\n")
		if len(foreignKeys) > 0 {
			allDefs += ",\n" + strings.Join(foreignKeys, ",\n")
		}
		sb.WriteString(allDefs)
		sb.WriteString("\n);\n\n")
		for _, field := range entity.Fields {
			if field.IsIndexed {
				var colName string
				if field.Relation != "" {
					colName = toSnakeCase(field.Name) + "_id"
				} else {
					colName = toSnakeCase(field.Name)
				}
				idxName := fmt.Sprintf("idx_%s_%s", tableName, toSnakeCase(field.Name))
				if field.IsUnique {
					sb.WriteString(fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\"(\"%s\");\n", idxName, tableName, colName))
				} else {
					sb.WriteString(fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\"(\"%s\");\n", idxName, tableName, colName))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

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
	return sb.String(), nil
}

func toSnakeCase(s string) string {
	var r []rune
	for i, ch := range s {
		if i > 0 && ch >= 'A' && ch <= 'Z' {
			r = append(r, '_')
		}
		r = append(r, ch)
	}
	return strings.ToLower(string(r))
}

func sortEntities(entities []Entity) ([]Entity, error) {
	entMap := make(map[string]*Entity)
	for i := range entities {
		entMap[entities[i].Name] = &entities[i]
	}
	deps := make(map[string][]string)
	for _, ent := range entities {
		for _, field := range ent.Fields {
			if field.Relation != "" {
				deps[ent.Name] = append(deps[ent.Name], field.Relation)
			}
		}
	}
	var sorted []string
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var visit func(name string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("cyclic dependency detected at %s", name)
		}
		if !visited[name] {
			temp[name] = true
			for _, dep := range deps[name] {
				if _, ok := entMap[dep]; ok {
					if err := visit(dep); err != nil {
						return err
					}
				}
			}
			temp[name] = false
			visited[name] = true
			sorted = append(sorted, name)
		}
		return nil
	}
	for name := range entMap {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	sortedEntities := make([]Entity, 0, len(entities))
	for _, name := range sorted {
		sortedEntities = append(sortedEntities, *entMap[name])
	}
	return sortedEntities, nil
}
