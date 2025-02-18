package core

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	db "github.com/u2u-labs/layerg-crawler/db/graphqldb"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

// QueryResolver represents the query resolver
type QueryResolver struct {
	Schema *Schema
}

// ExtractRequestedFields inspects p.Info to extract the list of requested field names.
func ExtractRequestedFields(info graphql.ResolveInfo) []string {
	var fields []string
	for _, f := range info.FieldASTs {
		if f.SelectionSet != nil {
			for _, sel := range f.SelectionSet.Selections {
				if field, ok := sel.(*ast.Field); ok {
					// Convert field name to snake_case
					fieldName := toSnakeCase(field.Name.Value)
					fields = append(fields, fieldName)
				}
			}
		}
	}
	return fields
}

// toSnakeCase converts a camelCase string to snake_case
func toSnakeCase(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += strings.ToLower(string(r))
	}
	return result
}

// toCamelCase converts a snake_case string to camelCase
func toCamelCase(s string) string {
	var result string
	capitalize := false
	for i, r := range s {
		if r == '_' {
			capitalize = true
			continue
		}
		if i == 0 {
			result += strings.ToLower(string(r))
		} else if capitalize {
			result += strings.ToUpper(string(r))
			capitalize = false
		} else {
			result += string(r)
		}
	}
	return result
}

// deriveTableName returns the table name based on the type name.
// If the lowercased typeName already ends in "s", we assume it's plural and use it as is.
// Otherwise, we append an "s".
func deriveTableName(typeName string) string {
	return toSnakeCase(typeName)
}

// isDateField checks if a field is a date type by looking at the GraphQL schema
func (r *QueryResolver) isDateField(typeName, fieldName string) bool {
	// Get the type definition from schema
	typeObj, exists := r.Schema.Types[typeName]
	if !exists {
		return false
	}

	// Look for the field in the type definition
	for _, field := range typeObj.Fields {
		if field.Name.Value == toCamelCase(fieldName) {
			// Check if it's a custom scalar type
			if field.Type.GetKind() == kinds.Named {
				namedType := field.Type.(*ast.Named)
				// Check if the type is one of our date scalars
				switch namedType.Name.Value {
				case "DateTime", "Date", "Time", "Timestamp":
					return true
				}
			}
		}
	}
	return false
}

// ResolveSingle builds a dynamic SQL query for a single record.
func (r *QueryResolver) ResolveSingle(typeName string, p graphql.ResolveParams) (interface{}, error) {
	requested := ExtractRequestedFields(p.Info)
	if len(requested) == 0 {
		requested = []string{"id"}
	}

	tableName := deriveTableName(typeName)

	// Create field mapping
	fieldToIndex := make(map[string]int)
	currentIndex := 0

	// Handle relationships
	joins := []string{}
	selectFields := []string{}

	// First collect all direct fields
	for _, field := range requested {
		switch field {
		case "posts":
			continue // Skip adding to selectFields
		case "author":
			joins = append(joins, `LEFT JOIN "user" ON "post"."author_id" = "user"."id"`)
			selectFields = append(selectFields, `"user"."id" as author_id, "user"."name" as author_name`)
			fieldToIndex[field] = currentIndex
			currentIndex++
		case "profile":
			joins = append(joins, `LEFT JOIN "user_profile" ON "user"."profile_id" = "user_profile"."id"`)
			selectFields = append(selectFields, fmt.Sprintf(`"user_profile"."id" as profile_id, "user_profile"."bio", "user_profile"."avatar_url"`))
			fieldToIndex[field] = currentIndex
			currentIndex++
		default:
			selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, field))
			fieldToIndex[field] = currentIndex
			currentIndex++
		}
	}

	// If no fields were added (only relationships requested), add id
	if len(selectFields) == 0 {
		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
		fieldToIndex["id"] = currentIndex
	}

	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s WHERE "%s"."id" = $1`,
		strings.Join(selectFields, ","),
		tableName,
		strings.Join(joins, " "),
		tableName)

	log.Println("SQL Query:", query)
	row := db.DB.QueryRow(query, p.Args["id"])
	values := make([]interface{}, len(selectFields))
	for i := range values {
		// Check if the current field is a date field
		var field string
		for f, idx := range fieldToIndex {
			if idx == i {
				field = f
				break
			}
		}

		if r.isDateField(typeName, field) {
			values[i] = &sql.NullTime{}
		} else {
			values[i] = &sql.NullString{}
		}
	}
	err := row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	result := map[string]interface{}{}

	// Map values using fieldToIndex and convert back to camelCase
	for field, idx := range fieldToIndex {
		if r.isDateField(typeName, field) {
			val := values[idx].(*sql.NullTime)
			if val.Valid {
				camelField := toCamelCase(field)
				result[camelField] = val.Time
			}
		} else {
			val := values[idx].(*sql.NullString)
			if val.Valid {
				camelField := toCamelCase(field)
				result[camelField] = val.String
			}
		}
	}

	// Handle nested posts
	if result != nil {
		for _, field := range requested {
			if field == "posts" && typeName == "User" {
				postsQuery := `SELECT id, title, content, published_date FROM "post" WHERE "author_id" = $1`
				rows, err := db.DB.Query(postsQuery, p.Args["id"])
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var posts []map[string]interface{}
				for rows.Next() {
					var id, title, content sql.NullString
					var publishedDate sql.NullTime
					err := rows.Scan(&id, &title, &content, &publishedDate)
					if err != nil {
						return nil, err
					}
					post := map[string]interface{}{
						"id":            id.String,
						"title":         title.String,
						"content":       content.String,
						"publishedDate": publishedDate.Time,
					}
					posts = append(posts, post)
				}
				result["posts"] = posts
			}
		}
	}

	return result, nil
}

// ResolveMultiple builds a dynamic SQL query for multiple records.
func (r *QueryResolver) ResolveMultiple(typeName string, p graphql.ResolveParams) (interface{}, error) {
	requested := ExtractRequestedFields(p.Info)
	if len(requested) == 0 {
		requested = []string{"id"}
	}

	tableName := deriveTableName(typeName)

	// Handle relationships
	joins := []string{}
	selectFields := []string{}
	hasNestedFields := false
	fieldToIndex := make(map[string]int) // Track field positions
	currentIndex := 0

	// Always include ID when posts are requested
	needsId := false
	for _, field := range requested {
		if field == "posts" {
			needsId = true
			hasNestedFields = true
			break
		}
	}

	// If we need ID, add it first
	if needsId {
		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
		fieldToIndex["id"] = currentIndex
		currentIndex++
	}

	// Add other requested fields
	for _, field := range requested {
		switch field {
		case "posts":
			continue
		case "author":
			joins = append(joins, `LEFT JOIN "user" ON "post"."author_id" = "user"."id"`)
			selectFields = append(selectFields, `"user"."id" as author_id, "user"."name" as author_name`)
			fieldToIndex[field] = currentIndex
			currentIndex++
		case "profile":
			joins = append(joins, `LEFT JOIN "user_profile" ON "user"."profile_id" = "user_profile"."id"`)
			selectFields = append(selectFields, `"user_profile"."id" as profile_id, "user_profile"."bio", "user_profile"."avatar_url"`)
			fieldToIndex[field] = currentIndex
			currentIndex++
		default:
			if field != "id" || !needsId { // Skip id if already added
				selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, field))
				fieldToIndex[field] = currentIndex
				currentIndex++
			}
		}
	}

	// If no fields were added, add id
	if len(selectFields) == 0 {
		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
		fieldToIndex["id"] = 0
	}

	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s`,
		strings.Join(selectFields, ","),
		tableName,
		strings.Join(joins, " "))

	// Handle pagination
	if page, ok := p.Args["page"].(int); ok {
		limit := 10 // default limit
		if lim, ok := p.Args["limit"].(int); ok {
			limit = lim
		}
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	log.Println("SQL Query:", query)
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(selectFields))
		for i := range values {
			// Check if the current field is a date field
			var field string
			for f, idx := range fieldToIndex {
				if idx == i {
					field = f
					break
				}
			}

			if r.isDateField(typeName, field) {
				values[i] = &sql.NullTime{}
			} else {
				values[i] = &sql.NullString{}
			}
		}
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		record := map[string]interface{}{}

		// Map values to fields using fieldToIndex and convert back to camelCase
		for field, idx := range fieldToIndex {
			if r.isDateField(typeName, field) {
				val := values[idx].(*sql.NullTime)
				if val.Valid {
					camelField := toCamelCase(field)
					record[camelField] = val.Time
				}
			} else {
				val := values[idx].(*sql.NullString)
				if val.Valid {
					camelField := toCamelCase(field)
					record[camelField] = val.String
				}
			}
		}

		// Handle nested posts if requested
		if hasNestedFields {
			userId, ok := record["id"].(string)
			if ok {
				postsQuery := `SELECT id, content FROM "post" WHERE "author_id" = $1`
				postRows, err := db.DB.Query(postsQuery, userId)
				if err != nil {
					return nil, err
				}
				defer postRows.Close()

				var posts []map[string]interface{}
				for postRows.Next() {
					var id, content sql.NullString
					err := postRows.Scan(&id, &content)
					if err != nil {
						return nil, err
					}
					post := map[string]interface{}{}
					if id.Valid {
						post["id"] = id.String
					}
					if content.Valid {
						post["content"] = content.String
					}
					posts = append(posts, post)
				}
				record["posts"] = posts
			}
		}

		results = append(results, record)
	}
	return results, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
