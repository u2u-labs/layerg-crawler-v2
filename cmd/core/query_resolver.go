package core

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
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
	// Base query.
	query := fmt.Sprintf(`SELECT DISTINCT %s FROM "%s" %s`,
		strings.Join(selectFields, ", "),
		tableName,
		strings.Join(joins, " "))
	// Build WHERE clause.
	args := []interface{}{}
	if whereRaw, ok := p.Args["where"]; ok && whereRaw != nil {
		if whereMap, ok := whereRaw.(map[string]interface{}); ok {
			conditions := []string{}
			for field, filter := range whereMap {
				if filterMap, ok := filter.(map[string]interface{}); ok {
					columnName := fmt.Sprintf(`"%s"."%s"`, tableName, toSnakeCase(field))
					for op, val := range filterMap {
						if val == nil {
							continue
						}
						var converted interface{} = val
						switch v := val.(type) {
						case big.Int:
							converted = v.String()
						case *big.Int:
							converted = v.String()
						}
						var sqlOp string
						switch op {
						case "gte":
							sqlOp = ">="
						case "gt":
							sqlOp = ">"
						case "eq":
							sqlOp = "="
						case "lt":
							sqlOp = "<"
						case "lte":
							sqlOp = "<="
						default:
							continue
						}
						conditions = append(conditions, fmt.Sprintf("%s %s $%d", columnName, sqlOp, len(args)+1))
						args = append(args, converted)
					}
				}
			}
			if len(conditions) > 0 {
				query += " WHERE " + strings.Join(conditions, " AND ")
			}
		}
	}
	// Append ORDER BY, LIMIT, and OFFSET.
	if order, ok := p.Args["order"].(string); ok && order != "" {
		query += " ORDER BY " + order
	}
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
	rows, err := db.DB.Query(query, args...)
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
		// Handle derived (inverse) fields.
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
			var nestedAliases []string
			if fieldAST.SelectionSet != nil {
				for _, sel := range fieldAST.SelectionSet.Selections {
					if sf, ok := sel.(*ast.Field); ok {
						fieldNameNested := sf.Name.Value
						nestedFieldDef := r.getFieldDefinition(relatedType, toCamelCase(fieldNameNested))
						var colName string
						if nestedFieldDef != nil && !isScalar(nestedFieldDef.Type) {
							colName = toSnakeCase(fieldNameNested) + "_id"
						} else {
							colName = toSnakeCase(fieldNameNested)
						}
						alias := toCamelCase(fieldNameNested)
						nestedSelects = append(nestedSelects, fmt.Sprintf(`"%s"."%s" as "%s"`, relatedTable, colName, alias))
						nestedAliases = append(nestedAliases, alias)
					}
				}
			}
			if len(nestedSelects) == 0 {
				nestedSelects = append(nestedSelects, fmt.Sprintf(`"%s"."id" as "id"`, relatedTable))
				nestedAliases = append(nestedAliases, "id")
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
			for nestedRows.Next() {
				vals := make([]interface{}, len(nestedAliases))
				for i := range vals {
					vals[i] = new(sql.NullString)
				}
				if err := nestedRows.Scan(vals...); err != nil {
					nestedRows.Close()
					return nil, err
				}
				rec := make(map[string]interface{})
				for i, alias := range nestedAliases {
					ns := vals[i].(*sql.NullString)
					var v interface{}
					if ns.Valid {
						v = ns.String
					}
					nestedFieldDef := r.getFieldDefinition(relatedType, alias)
					if nestedFieldDef != nil && !isScalar(nestedFieldDef.Type) {
						rec[alias] = map[string]interface{}{"id": v}
					} else {
						rec[alias] = v
					}
				}
				nestedResults = append(nestedResults, rec)
			}
			nestedRows.Close()
			record[toCamelCase(fieldName)] = nestedResults
		}
		// Post-process join fields: group flat join columns into nested objects.
		for _, fieldDef := range r.Schema.Types[typeName].Fields {
			if !isScalar(fieldDef.Type) && !hasDirective(fieldDef, "derivedFrom") {
				prefix := toCamelCase(toSnakeCase(fieldDef.Name.Value))
				nested := make(map[string]interface{})
				for k, v := range record {
					if k != prefix && strings.HasPrefix(k, prefix) {
						subKey := k[len(prefix):]
						if len(subKey) > 0 {
							subKey = strings.ToLower(subKey[:1]) + subKey[1:]
							nested[subKey] = v
							delete(record, k)
						}
					}
				}
				if len(nested) > 0 {
					record[prefix] = nested
				}
			}
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
		strings.Join(selectFields, ","),
		tableName,
		strings.Join(joins, " "),
		tableName)
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
	// Handle derived (inverse) fields.
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
				result[toCamelCase(fname)] = nestedResults
			}
		}
	}
	// Post-process join fields: group flat join columns into nested objects.
	for _, fieldDef := range r.Schema.Types[typeName].Fields {
		if !isScalar(fieldDef.Type) && !hasDirective(fieldDef, "derivedFrom") {
			prefix := toCamelCase(toSnakeCase(fieldDef.Name.Value))
			nested := make(map[string]interface{})
			for k, v := range result {
				if k != prefix && strings.HasPrefix(k, prefix) {
					subKey := k[len(prefix):]
					if len(subKey) > 0 {
						subKey = strings.ToLower(subKey[:1]) + subKey[1:]
						nested[subKey] = v
						delete(result, k)
					}
				}
			}
			if len(nested) > 0 {
				result[prefix] = nested
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
