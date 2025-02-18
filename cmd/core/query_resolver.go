// package core

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"strings"

// 	db "github.com/u2u-labs/layerg-crawler/db/graphqldb"

// 	"github.com/graphql-go/graphql"
// 	"github.com/graphql-go/graphql/language/ast"
// 	"github.com/graphql-go/graphql/language/kinds"
// )

// // QueryResolver represents the query resolver
// type QueryResolver struct {
// 	Schema *Schema
// }

// // ExtractRequestedFields inspects p.Info to extract the list of requested field names.
// func ExtractRequestedFields(info graphql.ResolveInfo) []string {
// 	var fields []string
// 	for _, f := range info.FieldASTs {
// 		if f.SelectionSet != nil {
// 			for _, sel := range f.SelectionSet.Selections {
// 				if field, ok := sel.(*ast.Field); ok {
// 					// Convert field name to snake_case
// 					fieldName := toSnakeCase(field.Name.Value)
// 					fields = append(fields, fieldName)
// 				}
// 			}
// 		}
// 	}
// 	return fields
// }

// // toSnakeCase converts a camelCase string to snake_case
// func toSnakeCase(s string) string {
// 	var result string
// 	for i, r := range s {
// 		if i > 0 && r >= 'A' && r <= 'Z' {
// 			result += "_"
// 		}
// 		result += strings.ToLower(string(r))
// 	}
// 	return result
// }

// // toCamelCase converts a snake_case string to camelCase
// func toCamelCase(s string) string {
// 	var result string
// 	capitalize := false
// 	for i, r := range s {
// 		if r == '_' {
// 			capitalize = true
// 			continue
// 		}
// 		if i == 0 {
// 			result += strings.ToLower(string(r))
// 		} else if capitalize {
// 			result += strings.ToUpper(string(r))
// 			capitalize = false
// 		} else {
// 			result += string(r)
// 		}
// 	}
// 	return result
// }

// // deriveTableName returns the table name based on the type name.
// // If the lowercased typeName already ends in "s", we assume it's plural and use it as is.
// // Otherwise, we append an "s".
// func deriveTableName(typeName string) string {
// 	return toSnakeCase(typeName)
// }

// // isDateField checks if a field is a date type by looking at the GraphQL schema
// func (r *QueryResolver) isDateField(typeName, fieldName string) bool {
// 	// Get the type definition from schema
// 	typeObj, exists := r.Schema.Types[typeName]
// 	if !exists {
// 		return false
// 	}

// 	// Look for the field in the type definition
// 	for _, field := range typeObj.Fields {
// 		if field.Name.Value == toCamelCase(fieldName) {
// 			// Check if it's a custom scalar type
// 			if field.Type.GetKind() == kinds.Named {
// 				namedType := field.Type.(*ast.Named)
// 				// Check if the type is one of our date scalars
// 				switch namedType.Name.Value {
// 				case "DateTime", "Date", "Time", "Timestamp":
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// // ResolveSingle builds a dynamic SQL query for a single record.
// func (r *QueryResolver) ResolveSingle(typeName string, p graphql.ResolveParams) (interface{}, error) {
// 	requested := ExtractRequestedFields(p.Info)
// 	if len(requested) == 0 {
// 		requested = []string{"id"}
// 	}

// 	tableName := deriveTableName(typeName)

// 	// Create field mapping
// 	fieldToIndex := make(map[string]int)
// 	currentIndex := 0

// 	// Handle relationships
// 	joins := []string{}
// 	selectFields := []string{}

// 	// First collect all direct fields
// 	for _, field := range requested {
// 		switch field {
// 		case "posts":
// 			continue // Skip adding to selectFields
// 		case "author":
// 			joins = append(joins, `LEFT JOIN "user" ON "post"."author_id" = "user"."id"`)
// 			selectFields = append(selectFields, `"user"."id" as author_id, "user"."name" as author_name`)
// 			fieldToIndex[field] = currentIndex
// 			currentIndex++
// 		case "profile":
// 			joins = append(joins, `LEFT JOIN "user_profile" ON "user"."profile_id" = "user_profile"."id"`)
// 			selectFields = append(selectFields, fmt.Sprintf(`"user_profile"."id" as profile_id, "user_profile"."bio", "user_profile"."avatar_url"`))
// 			fieldToIndex[field] = currentIndex
// 			currentIndex++
// 		default:
// 			selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, field))
// 			fieldToIndex[field] = currentIndex
// 			currentIndex++
// 		}
// 	}

// 	// If no fields were added (only relationships requested), add id
// 	if len(selectFields) == 0 {
// 		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
// 		fieldToIndex["id"] = currentIndex
// 	}

// 	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s WHERE "%s"."id" = $1`,
// 		strings.Join(selectFields, ","),
// 		tableName,
// 		strings.Join(joins, " "),
// 		tableName)

// 	log.Println("SQL Query:", query)
// 	row := db.DB.QueryRow(query, p.Args["id"])
// 	values := make([]interface{}, len(selectFields))
// 	for i := range values {
// 		// Check if the current field is a date field
// 		var field string
// 		for f, idx := range fieldToIndex {
// 			if idx == i {
// 				field = f
// 				break
// 			}
// 		}

// 		if r.isDateField(typeName, field) {
// 			values[i] = &sql.NullTime{}
// 		} else {
// 			values[i] = &sql.NullString{}
// 		}
// 	}
// 	err := row.Scan(values...)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	result := map[string]interface{}{}

// 	// Map values using fieldToIndex and convert back to camelCase
// 	for field, idx := range fieldToIndex {
// 		if r.isDateField(typeName, field) {
// 			val := values[idx].(*sql.NullTime)
// 			if val.Valid {
// 				camelField := toCamelCase(field)
// 				result[camelField] = val.Time
// 			}
// 		} else {
// 			val := values[idx].(*sql.NullString)
// 			if val.Valid {
// 				camelField := toCamelCase(field)
// 				result[camelField] = val.String
// 			}
// 		}
// 	}

// 	// Handle nested posts
// 	if result != nil {
// 		for _, field := range requested {
// 			if field == "posts" && typeName == "User" {
// 				postsQuery := `SELECT id, title, content, published_date FROM "post" WHERE "author_id" = $1`
// 				rows, err := db.DB.Query(postsQuery, p.Args["id"])
// 				if err != nil {
// 					return nil, err
// 				}
// 				defer rows.Close()

// 				var posts []map[string]interface{}
// 				for rows.Next() {
// 					var id, title, content sql.NullString
// 					var publishedDate sql.NullTime
// 					err := rows.Scan(&id, &title, &content, &publishedDate)
// 					if err != nil {
// 						return nil, err
// 					}
// 					post := map[string]interface{}{
// 						"id":            id.String,
// 						"title":         title.String,
// 						"content":       content.String,
// 						"publishedDate": publishedDate.Time,
// 					}
// 					posts = append(posts, post)
// 				}
// 				result["posts"] = posts
// 			}
// 		}
// 	}

// 	return result, nil
// }

// // ResolveMultiple builds a dynamic SQL query for multiple records.
// func (r *QueryResolver) ResolveMultiple(typeName string, p graphql.ResolveParams) (interface{}, error) {
// 	requested := ExtractRequestedFields(p.Info)
// 	if len(requested) == 0 {
// 		requested = []string{"id"}
// 	}

// 	tableName := deriveTableName(typeName)

// 	// Handle relationships
// 	joins := []string{}
// 	selectFields := []string{}
// 	hasNestedFields := false
// 	fieldToIndex := make(map[string]int) // Track field positions
// 	currentIndex := 0

// 	// Always include ID when posts are requested
// 	needsId := false
// 	for _, field := range requested {
// 		if field == "posts" {
// 			needsId = true
// 			hasNestedFields = true
// 			break
// 		}
// 	}

// 	// If we need ID, add it first
// 	if needsId {
// 		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
// 		fieldToIndex["id"] = currentIndex
// 		currentIndex++
// 	}

// 	// Add other requested fields
// 	for _, field := range requested {
// 		switch field {
// 		case "posts":
// 			continue
// 		case "author":
// 			joins = append(joins, `LEFT JOIN "user" ON "post"."author_id" = "user"."id"`)
// 			selectFields = append(selectFields, `"user"."id" as author_id, "user"."name" as author_name`)
// 			fieldToIndex[field] = currentIndex
// 			currentIndex++
// 		case "profile":
// 			joins = append(joins, `LEFT JOIN "user_profile" ON "user"."profile_id" = "user_profile"."id"`)
// 			selectFields = append(selectFields, `"user_profile"."id" as profile_id, "user_profile"."bio", "user_profile"."avatar_url"`)
// 			fieldToIndex[field] = currentIndex
// 			currentIndex++
// 		default:
// 			if field != "id" || !needsId { // Skip id if already added
// 				selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, field))
// 				fieldToIndex[field] = currentIndex
// 				currentIndex++
// 			}
// 		}
// 	}

// 	// If no fields were added, add id
// 	if len(selectFields) == 0 {
// 		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
// 		fieldToIndex["id"] = 0
// 	}

// 	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s`,
// 		strings.Join(selectFields, ","),
// 		tableName,
// 		strings.Join(joins, " "))

// 	// Handle pagination
// 	if page, ok := p.Args["page"].(int); ok {
// 		limit := 10 // default limit
// 		if lim, ok := p.Args["limit"].(int); ok {
// 			limit = lim
// 		}
// 		offset := (page - 1) * limit
// 		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
// 	}

// 	log.Println("SQL Query:", query)
// 	rows, err := db.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var results []map[string]interface{}
// 	for rows.Next() {
// 		values := make([]interface{}, len(selectFields))
// 		for i := range values {
// 			// Check if the current field is a date field
// 			var field string
// 			for f, idx := range fieldToIndex {
// 				if idx == i {
// 					field = f
// 					break
// 				}
// 			}

// 			if r.isDateField(typeName, field) {
// 				values[i] = &sql.NullTime{}
// 			} else {
// 				values[i] = &sql.NullString{}
// 			}
// 		}
// 		err := rows.Scan(values...)
// 		if err != nil {
// 			return nil, err
// 		}

// 		record := map[string]interface{}{}

// 		// Map values to fields using fieldToIndex and convert back to camelCase
// 		for field, idx := range fieldToIndex {
// 			if r.isDateField(typeName, field) {
// 				val := values[idx].(*sql.NullTime)
// 				if val.Valid {
// 					camelField := toCamelCase(field)
// 					record[camelField] = val.Time
// 				}
// 			} else {
// 				val := values[idx].(*sql.NullString)
// 				if val.Valid {
// 					camelField := toCamelCase(field)
// 					record[camelField] = val.String
// 				}
// 			}
// 		}

// 		// Handle nested posts if requested
// 		if hasNestedFields {
// 			userId, ok := record["id"].(string)
// 			if ok {
// 				postsQuery := `SELECT id, content FROM "post" WHERE "author_id" = $1`
// 				postRows, err := db.DB.Query(postsQuery, userId)
// 				if err != nil {
// 					return nil, err
// 				}
// 				defer postRows.Close()

// 				var posts []map[string]interface{}
// 				for postRows.Next() {
// 					var id, content sql.NullString
// 					err := postRows.Scan(&id, &content)
// 					if err != nil {
// 						return nil, err
// 					}
// 					post := map[string]interface{}{}
// 					if id.Valid {
// 						post["id"] = id.String
// 					}
// 					if content.Valid {
// 						post["content"] = content.String
// 					}
// 					posts = append(posts, post)
// 				}
// 				record["posts"] = posts
// 			}
// 		}

// 		results = append(results, record)
// 	}
// 	return results, nil
// }

// // Helper function to check if a slice contains a string
// func contains(slice []string, str string) bool {
// 	for _, v := range slice {
// 		if v == str {
// 			return true
// 		}
// 	}
// 	return false
// }

package core

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	db "github.com/u2u-labs/layerg-crawler/db/graphqldb"
)

// getFieldDefinition looks up the field definition in the schema.
type QueryResolver struct {
	Schema *Schema
}

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

func deriveTableName(typeName string) string {
	return toSnakeCase(typeName)
}

func unwrapType(t ast.Type) *ast.Named {
	switch typ := t.(type) {
	case *ast.NonNull:
		return unwrapType(typ.Type)
	case *ast.List:
		return unwrapType(typ.Type)
	case *ast.Named:
		return typ
	}
	return nil
}

func isScalar(t ast.Type) bool {
	named := unwrapType(t)
	if named == nil {
		return false
	}
	switch named.Name.Value {
	case "ID", "String", "BigInt", "Int", "Float", "Boolean", "DateTime", "Date", "Time", "Timestamp":
		return true
	}
	return false
}

func (r *QueryResolver) isDateField(typeName, fieldName string) bool {
	typeObj, exists := r.Schema.Types[typeName]
	if !exists {
		return false
	}
	// Convert fieldName to camelCase to match the schema definition.
	camelField := toCamelCase(fieldName)
	for _, field := range typeObj.Fields {
		if field.Name.Value == camelField {
			namedType := unwrapType(field.Type)
			if namedType == nil {
				return false
			}
			switch namedType.Name.Value {
			case "DateTime", "Date", "Time", "Timestamp":
				return true
			}
		}
	}
	return false
}

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

func (r *QueryResolver) getFieldDefinition(typeName, fieldName string) *ast.FieldDefinition {
	if typeDef, ok := r.Schema.Types[typeName]; ok {
		for _, field := range typeDef.Fields {
			if field.Name.Value == toCamelCase(fieldName) || field.Name.Value == fieldName {
				return field
			}
		}
	}
	return nil
}

// hasDirective returns true if the field has a directive with the given name.
func hasDirective(field *ast.FieldDefinition, name string) bool {
	for _, d := range field.Directives {
		if d.Name.Value == name {
			return true
		}
	}
	return false
}

// getNamedType extracts the underlying named type from a field type.
func getNamedType(t ast.Node) string {
	switch typ := t.(type) {
	case *ast.Named:
		return typ.Name.Value
	case *ast.List:
		return getNamedType(typ.Type)
	case *ast.NonNull:
		return getNamedType(typ.Type)
	}
	return ""
}

func (r *QueryResolver) buildSelectAndJoinsWithDerived(typeName string, fields []*ast.Field) (selectFields []string, joins []string, fieldToIndex map[string]int, derivedFields map[string]*ast.Field) {
	tableName := deriveTableName(typeName)
	fieldToIndex = make(map[string]int)
	derivedFields = make(map[string]*ast.Field)
	currentIndex := 0
	for _, f := range fields {
		fname := toSnakeCase(f.Name.Value)
		fieldDef := r.getFieldDefinition(typeName, fname)
		if fieldDef == nil {
			// Skip if the schema doesn’t define this field.
			continue
		}
		if hasDirective(fieldDef, "derivedFrom") {
			// Do not add a column or join; resolve later.
			derivedFields[fname] = f
			continue
		}
		// If the field is scalar, simply select it.
		if isScalar(fieldDef.Type) {
			selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, fname))
			fieldToIndex[fname] = currentIndex
			currentIndex++
		} else {
			// For non–derived object relations (assumed stored as a FK)
			relatedType := getNamedType(fieldDef.Type)
			joinTable := deriveTableName(relatedType)
			// Convention: the FK column in the parent table is named fieldName + "_id"
			joinClause := fmt.Sprintf(`LEFT JOIN "%s" ON "%s"."%s_id" = "%s"."id"`, joinTable, tableName, fname, joinTable)
			joins = append(joins, joinClause)
			if f.SelectionSet != nil {
				for _, subSel := range f.SelectionSet.Selections {
					if sf, ok := subSel.(*ast.Field); ok {
						subName := toSnakeCase(sf.Name.Value)
						alias := fmt.Sprintf("%s_%s", fname, subName)
						selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s" as %s`, joinTable, subName, alias))
						fieldToIndex[alias] = currentIndex
						currentIndex++
					}
				}
			}
		}
	}
	return
}

func (r *QueryResolver) ResolveMultiple(typeName string, p graphql.ResolveParams) (interface{}, error) {
	tableName := deriveTableName(typeName)
	var selections []*ast.Field
	if len(p.Info.FieldASTs) > 0 {
		topField := p.Info.FieldASTs[0]
		if topField.SelectionSet != nil {
			for _, sel := range topField.SelectionSet.Selections {
				if sf, ok := sel.(*ast.Field); ok {
					selections = append(selections, sf)
				}
			}
		}
	}
	if len(selections) == 0 {
		selections = append(selections, &ast.Field{
			Name: &ast.Name{Value: "id"},
		})
	}
	// Build SELECT columns, JOIN clauses, and capture derived fields.
	selectFields, joins, fieldToIndex, derivedFields := r.buildSelectAndJoinsWithDerived(typeName, selections)
	// Ensure the primary key "id" is always selected.
	if _, ok := fieldToIndex["id"]; !ok {
		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
		fieldToIndex["id"] = len(selectFields) - 1
	}
	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s`,
		strings.Join(selectFields, ", "),
		tableName,
		strings.Join(joins, " "))
	// Add ordering if provided.
	if order, ok := p.Args["order"].(string); ok && order != "" {
		query += " ORDER BY " + order
	}
	// Add LIMIT (and OFFSET if page is provided) if a limit argument is passed.
	if lim, ok := p.Args["limit"].(int); ok {
		offset := 0
		if page, ok := p.Args["page"].(int); ok && page > 1 {
			offset = (page - 1) * lim
		}
		query += fmt.Sprintf(" LIMIT %d", lim)
		if offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", offset)
		}
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
			var fieldName string
			for f, idx := range fieldToIndex {
				if idx == i {
					fieldName = f
					break
				}
			}
			if r.isDateField(typeName, fieldName) {
				values[i] = new(sql.NullTime)
			} else {
				values[i] = new(sql.NullString)
			}
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		record := make(map[string]interface{})
		for field, idx := range fieldToIndex {
			camelField := toCamelCase(field)
			if r.isDateField(typeName, field) {
				val := values[idx].(*sql.NullTime)
				if val.Valid {
					record[camelField] = val.Time
				}
			} else {
				val := values[idx].(*sql.NullString)
				if val.Valid {
					record[camelField] = val.String
				}
			}
		}
		// Resolve derived fields (like "items" on User) via nested queries.
		for fieldName, fieldAST := range derivedFields {
			fieldDef := r.getFieldDefinition(typeName, fieldName)
			if fieldDef == nil {
				continue
			}
			var derivedFieldName string
			for _, d := range fieldDef.Directives {
				if d.Name.Value == "derivedFrom" && len(d.Arguments) > 0 {
					if argVal, ok := d.Arguments[0].Value.GetValue().(string); ok {
						derivedFieldName = argVal
					}
				}
			}
			if derivedFieldName == "" {
				continue
			}
			relatedType := getNamedType(fieldDef.Type)
			relatedTable := deriveTableName(relatedType)
			var nestedSelects []string
			if fieldAST.SelectionSet != nil {
				for _, sel := range fieldAST.SelectionSet.Selections {
					if sf, ok := sel.(*ast.Field); ok {
						nestedSelects = append(nestedSelects, fmt.Sprintf(`"%s"."%s"`, relatedTable, toSnakeCase(sf.Name.Value)))
					}
				}
			}
			if len(nestedSelects) == 0 {
				nestedSelects = append(nestedSelects, fmt.Sprintf(`"%s"."id"`, relatedTable))
			}
			joinColumn := fmt.Sprintf(`"%s_id"`, derivedFieldName)
			nestedQuery := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" WHERE %s = $1`,
				strings.Join(nestedSelects, ", "),
				relatedTable,
				joinColumn)
			parentID, ok := record["id"].(string)
			if !ok {
				continue
			}
			nestedRows, err := db.DB.Query(nestedQuery, parentID)
			if err != nil {
				return nil, err
			}
			var nestedResults []map[string]interface{}
			cols, _ := nestedRows.Columns()
			for nestedRows.Next() {
				vals := make([]interface{}, len(cols))
				for i := range vals {
					vals[i] = new(sql.NullString)
				}
				if err := nestedRows.Scan(vals...); err != nil {
					nestedRows.Close()
					return nil, err
				}
				rec := make(map[string]interface{})
				for i, col := range cols {
					ns := vals[i].(*sql.NullString)
					if ns.Valid {
						rec[toCamelCase(col)] = ns.String
					}
				}
				nestedResults = append(nestedResults, rec)
			}
			nestedRows.Close()
			record[toCamelCase(fieldName)] = nestedResults
		}
		results = append(results, record)
	}
	return results, nil
}

// ResolveSingle builds a dynamic SQL query based on the schema.
func (r *QueryResolver) ResolveSingle(typeName string, p graphql.ResolveParams) (interface{}, error) {
	tableName := deriveTableName(typeName)
	selectFields, joins, fieldToIndex := r.buildSelectAndJoins(typeName, p.Info.FieldASTs)
	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s WHERE "%s"."id" = $1`,
		strings.Join(selectFields, ","), tableName, strings.Join(joins, " "), tableName)
	log.Println("SQL Query:", query)
	row := db.DB.QueryRow(query, p.Args["id"])
	values := make([]interface{}, len(selectFields))
	for i := range values {
		var field string
		for k, idx := range fieldToIndex {
			if idx == i {
				field = k
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
	for field, idx := range fieldToIndex {
		if r.isDateField(typeName, field) {
			val := values[idx].(*sql.NullTime)
			if val.Valid {
				result[toCamelCase(field)] = val.Time
			}
		} else {
			val := values[idx].(*sql.NullString)
			if val.Valid {
				result[toCamelCase(field)] = val.String
			}
		}
	}
	// Dynamically handle fields with @derivedFrom (inverse relationships)
	if _, exists := r.Schema.Types[typeName]; exists {
		for _, f := range p.Info.FieldASTs {
			fname := toSnakeCase(f.Name.Value)
			fieldDef := r.getFieldDefinition(typeName, fname)
			if fieldDef != nil && !isScalar(fieldDef.Type) && hasDirective(fieldDef, "derivedFrom") {
				derivedField := ""
				for _, d := range fieldDef.Directives {
					if d.Name.Value == "derivedFrom" && len(d.Arguments) > 0 {
						derivedField = d.Arguments[0].Value.GetValue().(string)
					}
				}
				relatedType := getNamedType(fieldDef.Type)
				joinTable := deriveTableName(relatedType)
				nestedQuery := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s" = $1`, joinTable, derivedField+"_id")
				nestedRows, err := db.DB.Query(nestedQuery, p.Args["id"])
				if err != nil {
					return nil, err
				}
				defer nestedRows.Close()
				var nestedResults []map[string]interface{}
				cols, _ := nestedRows.Columns()
				for nestedRows.Next() {
					vals := make([]interface{}, len(cols))
					for i := range vals {
						vals[i] = new(sql.NullString)
					}
					if err := nestedRows.Scan(vals...); err != nil {
						return nil, err
					}
					rec := make(map[string]interface{})
					for i, col := range cols {
						ns := vals[i].(*sql.NullString)
						if ns.Valid {
							rec[toCamelCase(col)] = ns.String
						}
					}
					nestedResults = append(nestedResults, rec)
				}
				result[toCamelCase(fname)] = nestedResults
			}
		}
	}
	return result, nil
}

func (r *QueryResolver) buildSelectAndJoins(typeName string, fields []*ast.Field) (selectFields []string, joins []string, fieldToIndex map[string]int) {
	tableName := deriveTableName(typeName)
	fieldToIndex = make(map[string]int)
	currentIndex := 0
	for _, f := range fields {
		fname := toSnakeCase(f.Name.Value)
		fieldDef := r.getFieldDefinition(typeName, fname)
		if fieldDef == nil || isScalar(fieldDef.Type) {
			selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s"`, tableName, fname))
			fieldToIndex[fname] = currentIndex
			currentIndex++
		} else {
			// For non-scalar fields, e.g. relationships
			if hasDirective(fieldDef, "derivedFrom") {
				continue // will be resolved separately
			}
			relatedType := getNamedType(fieldDef.Type)
			joinTable := deriveTableName(relatedType)
			joinClause := fmt.Sprintf(`LEFT JOIN "%s" ON "%s"."%s_id" = "%s"."id"`, joinTable, tableName, fname, joinTable)
			joins = append(joins, joinClause)
			if f.SelectionSet != nil {
				for _, subSel := range f.SelectionSet.Selections {
					if subField, ok := subSel.(*ast.Field); ok {
						subName := toSnakeCase(subField.Name.Value)
						alias := fmt.Sprintf("%s_%s", fname, subName)
						selectFields = append(selectFields, fmt.Sprintf(`"%s"."%s" as %s`, joinTable, subName, alias))
						fieldToIndex[alias] = currentIndex
						currentIndex++
					}
				}
			}
		}
	}
	if len(selectFields) == 0 {
		selectFields = append(selectFields, fmt.Sprintf(`"%s"."id"`, tableName))
		fieldToIndex["id"] = currentIndex
	}
	return
}
