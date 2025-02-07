package generator

import (
	"errors"
	"io/ioutil"
	"strings"
)

// Field represents a field in an entity.
type Field struct {
	Name      string
	Type      string
	IsNonNull bool
	IsIndexed bool
	IsUnique  bool
	Relation  string // non-scalar types are assumed to be relations
}

// Entity represents a GraphQL entity.
type Entity struct {
	Name   string
	Fields []Field
}

// ParseGraphQLSchema parses a GraphQL schema file and extracts entity definitions.
// (NOTE: this is a very naive parser meant only as an example.)
func ParseGraphQLSchema(path string) ([]Entity, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	schema := string(content)
	var entities []Entity
	lines := strings.Split(schema, "\n")
	var currentEntity *Entity
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type") && strings.Contains(line, "@entity") {
			// For example: "type User @entity {"
			parts := strings.Split(line, " ")
			if len(parts) < 2 {
				return nil, errors.New("invalid type definition")
			}
			entityName := parts[1]
			currentEntity = &Entity{Name: entityName}
		} else if currentEntity != nil {
			if strings.HasPrefix(line, "}") {
				entities = append(entities, *currentEntity)
				currentEntity = nil
			} else if line != "" {
				// Parse field line (e.g.: "name: String! @index(unique: true)")
				fieldParts := strings.Split(line, ":")
				if len(fieldParts) < 2 {
					continue
				}
				fieldName := strings.TrimSpace(fieldParts[0])
				rest := strings.TrimSpace(fieldParts[1])
				// Get type token and check for non-null.
				tokens := strings.Split(rest, " ")
				fieldType := strings.TrimSpace(tokens[0])
				isNonNull := strings.HasSuffix(fieldType, "!")
				if isNonNull {
					fieldType = strings.TrimSuffix(fieldType, "!")
				}
				field := Field{
					Name:      fieldName,
					Type:      fieldType,
					IsNonNull: isNonNull,
				}
				// Check for indexing directives.
				if strings.Contains(rest, "@index") {
					field.IsIndexed = true
					if strings.Contains(rest, "unique: true") {
						field.IsUnique = true
					}
				}
				// Mark as a relation if the type is not a common scalar.
				if fieldType != "ID" && fieldType != "String" && fieldType != "Boolean" && fieldType != "Date" {
					field.Relation = fieldType
				}
				currentEntity.Fields = append(currentEntity.Fields, field)
			}
		}
	}
	return entities, nil
}
