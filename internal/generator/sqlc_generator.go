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
		var insertCols []string
		var insertPhs []string
		var updateAssignments []string

		// Add id to insert columns
		insertCols = append(insertCols, "\"id\"")
		insertPhs = append(insertPhs, "$1")

		// For other fields, start placeholder count from 2
		phCount := 2
		for _, field := range entity.Fields {
			if strings.ToLower(field.Name) == "id" {
				continue
			}
			if field.Relation != "" {
				continue
			}

			colName := toSnakeCase(field.Name)
			insertCols = append(insertCols, fmt.Sprintf("\"%s\"", colName))
			insertPhs = append(insertPhs, fmt.Sprintf("$%d", phCount))
			updateAssignments = append(updateAssignments, fmt.Sprintf("\"%s\" = $%d", colName, phCount))
			phCount++
		}

		// Generate CRUD queries
		createQuery := fmt.Sprintf("-- name: Create%s :one\nINSERT INTO \"%s\" (%s) VALUES (%s) RETURNING *;\n\n",
			entity.Name, tableName,
			strings.Join(insertCols, ", "),
			strings.Join(insertPhs, ", "))

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
		var updateQuery string
		if len(updateAssignments) == 0 {
			// If there are no fields to update, skip generating the update query
			updateQuery = fmt.Sprintf("-- name: %s :exec\n-- Skip update query generation as there are no updateable fields\n\n",
				updateQueryName)
		} else {
			updateQuery = fmt.Sprintf("-- name: %s :one\nUPDATE \"%s\" SET %s WHERE id = $1 RETURNING *;\n\n",
				updateQueryName, tableName, strings.Join(updateAssignments, ", "))
		}
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
