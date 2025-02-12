package generator

import (
	"errors"
	"io/ioutil"
	"strings"
)

type Field struct {
	Name      string
	Type      string
	IsNonNull bool
	IsIndexed bool
	IsUnique  bool
	Relation  string // non-scalar types are assumed to be relations
}

type Entity struct {
	Name   string
	Fields []Field
}

func isScalar(t string) bool {
	l := strings.ToLower(t)
	return l == "id" || l == "string" || l == "boolean" || l == "date" || l == "bigint"
}

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
				fieldParts := strings.Split(line, ":")
				if len(fieldParts) < 2 {
					continue
				}
				fieldName := strings.TrimSpace(fieldParts[0])
				rest := strings.TrimSpace(fieldParts[1])
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
				if strings.Contains(rest, "@index") {
					field.IsIndexed = true
					if strings.Contains(rest, "unique: true") {
						field.IsUnique = true
					}
				}
				if !isScalar(fieldType) {
					field.Relation = fieldType
				}
				currentEntity.Fields = append(currentEntity.Fields, field)
			}
		}
	}
	return entities, nil
}
