package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func generateFullMigrationDown(entities []Entity) (string, error) {
	sortedEntities, err := sortEntities(entities)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	// Reverse the alterations done in the second pass (for array relations)
	for _, entity := range sortedEntities {
		for _, field := range entity.Fields {
			if field.Relation != "" && strings.HasPrefix(field.Type, "[") && !field.DerivedFrom {
				baseType := strings.Trim(strings.Trim(field.Type, "[]!"), "!")
				manyTableName := toSnakeCase(baseType)
				oneTableName := toSnakeCase(entity.Name)
				sb.WriteString(fmt.Sprintf(`DROP INDEX IF EXISTS "idx_%s_%s_id";`, manyTableName, oneTableName))
				sb.WriteString("\n")
				sb.WriteString(fmt.Sprintf(`ALTER TABLE "%s" DROP CONSTRAINT IF EXISTS "fk_%s_%s";`, manyTableName, manyTableName, oneTableName))
				sb.WriteString("\n")
				sb.WriteString(fmt.Sprintf(`ALTER TABLE "%s" DROP COLUMN IF EXISTS "%s_id";`, manyTableName, oneTableName))
				sb.WriteString("\n\n")
			}
		}
	}
	// Drop the tables created in the first pass in reverse order
	for i := len(sortedEntities) - 1; i >= 0; i-- {
		tableName := toSnakeCase(sortedEntities[i].Name)
		sb.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS \"%s\" CASCADE;\n", tableName))
	}
	return sb.String(), nil
}

func GenerateMigrationScripts(entities []Entity, outputDir string) error {
	migrationsDir := outputDir + "/migrations"

	// Only remove if directory exists
	if _, err := os.Stat(migrationsDir); err == nil {
		if err := os.RemoveAll(migrationsDir); err != nil {
			return fmt.Errorf("failed to remove existing migrations directory: %w", err)
		}
	}

	// Create migrations directory
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

	// generate down migration based on the same entities (or subset in diff migration)
	downSQL, err := generateFullMigrationDown(entities)
	if err != nil {
		return err
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

	if _, err := f.WriteString(downSQL); err != nil {
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

	var sb strings.Builder

	// First pass: Create all base tables
	for _, entity := range sortedEntities {
		tableName := toSnakeCase(entity.Name)
		sb.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))
		sb.WriteString("    \"id\" TEXT PRIMARY KEY,\n")

		var colDefs []string
		var foreignKeys []string

		for _, field := range entity.Fields {
			if strings.ToLower(field.Name) == "id" {
				continue
			}

			var colDef string
			if field.Relation != "" {
				// Skip array relations in first pass
				if strings.HasPrefix(field.Type, "[") {
					continue
				}
				// Handle one-to-one relations
				colDef = fmt.Sprintf("    \"%s_id\" TEXT", toSnakeCase(field.Name))
				foreignKeys = append(foreignKeys,
					fmt.Sprintf("    FOREIGN KEY (\"%s_id\") REFERENCES \"%s\"(\"id\") ON DELETE CASCADE",
						toSnakeCase(field.Name),
						toSnakeCase(field.Relation)))
			} else {
				// Handle scalar types
				colName := toSnakeCase(field.Name)
				colType := getSQLType(field.Type)
				colDef = fmt.Sprintf("    \"%s\" %s", colName, colType)
			}

			if field.IsNonNull {
				colDef += " NOT NULL"
			}

			colDefs = append(colDefs, colDef)
		}

		// Add created_at timestamp
		colDefs = append(colDefs, "    \"created_at\" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP")

		// Combine all column definitions
		sb.WriteString(strings.Join(colDefs, ",\n"))
		if len(foreignKeys) > 0 {
			sb.WriteString(",\n")
			sb.WriteString(strings.Join(foreignKeys, ",\n"))
		}
		sb.WriteString("\n);\n\n")

		// Create indexes
		for _, field := range entity.Fields {
			if field.IsIndexed {
				idxName := fmt.Sprintf("idx_%s_%s", tableName, toSnakeCase(field.Name))
				colName := toSnakeCase(field.Name)
				if field.Relation != "" && !strings.HasPrefix(field.Type, "[") {
					colName += "_id"
				}

				if field.IsUnique {
					sb.WriteString(fmt.Sprintf("CREATE UNIQUE INDEX \"%s\" ON \"%s\"(\"%s\");\n",
						idxName, tableName, colName))
				} else {
					sb.WriteString(fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\"(\"%s\");\n",
						idxName, tableName, colName))
				}
			}
		}

		// create composite indexes
		compositeIndex := entity.CompositeIndex
		if len(compositeIndex) > 0 {

			idxName := fmt.Sprintf("idx_%s", tableName)
			for i, col := range compositeIndex {
				colNames := make([]string, len(compositeIndex[i]))
				for j, c := range col {
					colNames[j] = fmt.Sprintf("\"%s\"", toSnakeCase(c))
				}

				sb.WriteString(fmt.Sprintf("CREATE INDEX \"%s_%d\" ON \"%s\"(%s);\n",
					idxName, i, tableName, strings.Join(colNames, ", ")))
			}

		}

		sb.WriteString("\n")
	}

	// Second pass: Add foreign keys for array relations
	for _, entity := range sortedEntities {
		for _, field := range entity.Fields {
			if field.Relation != "" && strings.HasPrefix(field.Type, "[") {
				baseType := strings.Trim(strings.Trim(field.Type, "[]!"), "!")
				manyTableName := toSnakeCase(baseType)
				oneTableName := toSnakeCase(entity.Name)

				// Skip adding column if it's a @derivedFrom field
				if !field.DerivedFrom {
					// Add foreign key column and constraint
					sb.WriteString(fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN IF NOT EXISTS "%s_id" TEXT;`,
						manyTableName, oneTableName))
					sb.WriteString("\n")

					sb.WriteString(fmt.Sprintf(`ALTER TABLE "%s" ADD CONSTRAINT IF NOT EXISTS "fk_%s_%s" 
						FOREIGN KEY ("%s_id") REFERENCES "%s"("id") ON DELETE CASCADE;`,
						manyTableName, manyTableName, oneTableName,
						oneTableName, oneTableName))
					sb.WriteString("\n")

					sb.WriteString(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS "idx_%s_%s_id" ON "%s"("%s_id");`,
						manyTableName, oneTableName, manyTableName, oneTableName))
					sb.WriteString("\n\n")
				}
			}
		}
	}

	return sb.String(), nil
}

func getSQLType(graphqlType string) string {
	switch strings.ToLower(graphqlType) {
	case "id":
		return "TEXT"
	case "string":
		return "TEXT"
	case "boolean":
		return "BOOLEAN"
	case "int":
		return "INTEGER"
	case "bigint":
		return "NUMERIC"
	case "float":
		return "DOUBLE PRECISION"
	case "date":
		return "TIMESTAMPTZ"
	default:
		return "TEXT"
	}
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
