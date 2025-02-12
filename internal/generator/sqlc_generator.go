package generator

import (
	"fmt"
	"os"
	"strings"
)

// GenerateSQLCQueries generates basic CRUD SQL queries for sqlc.
func GenerateSQLCQueries(entities []Entity, outputDir string) error {
	queriesDir := outputDir + "/queries"
	if err := os.MkdirAll(queriesDir, os.ModePerm); err != nil {
		return err
	}
	filePath := queriesDir + "/queries.sql"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var sb strings.Builder
	// For each entity, generate CRUD queries.
	for _, entity := range entities {
		tableName := toSnakeCase(entity.Name)
		// Build insert columns and placeholders (skip auto-generated primary key).
		var insertCols []string
		var insertPhs []string
		var updateAssignments []string
		// For update queries, $1 is reserved for id.
		for _, field := range entity.Fields {
			// Skip auto-generated primary key.
			if strings.ToLower(field.Name) == "id" && field.Relation == "" {
				continue
			}
			var colName string
			if field.Relation != "" {
				colName = toSnakeCase(field.Name) + "_id"
			} else {
				colName = toSnakeCase(field.Name)
			}
			insertCols = append(insertCols, fmt.Sprintf("\"%s\"", colName))
			// Placeholder index for insert is the current count+1.
			insertPhs = append(insertPhs, fmt.Sprintf("$%d", len(insertPhs)+1))
			// For update assignments, placeholders start at $2 (since $1 is id).
			updateAssignments = append(updateAssignments, fmt.Sprintf("\"%s\" = $%d", colName, len(updateAssignments)+2))
		}
		// Create query.
		createQueryName := fmt.Sprintf("Create%s", entity.Name)
		var createQuery string
		if len(insertCols) > 0 {
			createQuery = fmt.Sprintf("-- name: %s :one\nINSERT INTO \"%s\" (%s) VALUES (%s) RETURNING *;\n\n",
				createQueryName,
				tableName,
				strings.Join(insertCols, ", "),
				strings.Join(insertPhs, ", "))
		} else {
			// In case there are no insertable columns.
			createQuery = fmt.Sprintf("-- name: %s :one\nINSERT INTO \"%s\" DEFAULT VALUES RETURNING *;\n\n",
				createQueryName, tableName)
		}
		// Get query.
		getQueryName := fmt.Sprintf("Get%s", entity.Name)
		getQuery := fmt.Sprintf("-- name: %s :one\nSELECT * FROM \"%s\" WHERE id = $1;\n\n",
			getQueryName, tableName)
		// List query.
		listQueryName := fmt.Sprintf("List%s", entity.Name)
		listQuery := fmt.Sprintf("-- name: %s :many\nSELECT * FROM \"%s\";\n\n",
			listQueryName, tableName)
		// Update query.
		updateQueryName := fmt.Sprintf("Update%s", entity.Name)
		updateQuery := fmt.Sprintf("-- name: %s :one\nUPDATE \"%s\" SET %s WHERE id = $1 RETURNING *;\n\n",
			updateQueryName, tableName, strings.Join(updateAssignments, ", "))
		// Delete query.
		deleteQueryName := fmt.Sprintf("Delete%s", entity.Name)
		deleteQuery := fmt.Sprintf("-- name: %s :exec\nDELETE FROM \"%s\" WHERE id = $1;\n\n",
			deleteQueryName, tableName)
		sb.WriteString(createQuery)
		sb.WriteString(getQuery)
		sb.WriteString(listQuery)
		sb.WriteString(updateQuery)
		sb.WriteString(deleteQuery)
	}
	_, err = f.WriteString(sb.String())
	return err
}
