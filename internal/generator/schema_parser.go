package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Field struct {
	Name        string
	Type        string
	IsNonNull   bool
	IsIndexed   bool
	IsUnique    bool
	IsList      bool
	DerivedFrom bool
	Relation    string // non-scalar types are assumed to be relations
}

type Entity struct {
	Name           string
	Fields         []Field
	CompositeIndex [][]string
}

type RAWSchema struct {
	Name   string
	Header string
	Fields []string
}

func isScalar(t string) bool {
	l := strings.ToLower(t)
	return l == "id" || l == "string" || l == "boolean" || l == "date" || l == "bigint"
}

func ParseGraphQLSchema(path string) ([]Entity, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading schema file:", err)
		return nil, err
	}

	schema := string(content)

	// Split the schema into lines
	lines := strings.Split(schema, "\n")

	var rawSchemas []RAWSchema
	var entities []Entity
	var currentSchema RAWSchema

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "type") {
			if currentSchema.Name != "" {
				rawSchemas = append(rawSchemas, currentSchema)
			}
			currentSchema = RAWSchema{Name: strings.Fields(trimmedLine)[1]}
			currentSchema.Header = trimmedLine
		} else if strings.Contains(trimmedLine, "@entity") {
			currentSchema.Header += " " + trimmedLine
		} else if strings.Contains(trimmedLine, "@compositeIndexes") {
			currentSchema.Header += " " + trimmedLine
		} else if trimmedLine == "}" {
			rawSchemas = append(rawSchemas, currentSchema)
			currentSchema = RAWSchema{}
		} else if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
			currentSchema.Fields = append(currentSchema.Fields, trimmedLine)
		}
	}
	if currentSchema.Name != "" {
		rawSchemas = append(rawSchemas, currentSchema)
	}

	var currentEntity *Entity
	for _, rawSchema := range rawSchemas {

		// handle header
		entityName := strings.Fields(rawSchema.Header)[1]
		currentEntity = &Entity{Name: entityName}

		// Extract composite index
		if strings.Contains(rawSchema.Header, "@compositeIndexes") {
			start := strings.Index(rawSchema.Header, "(")
			end := strings.LastIndex(rawSchema.Header, ")")
			if start != -1 && end != -1 && start < end {
				indexContent := rawSchema.Header[start+1 : end]
				if strings.HasPrefix(indexContent, "fields:") {
					fieldsStart := strings.Index(indexContent, "[")
					fieldsEnd := strings.LastIndex(indexContent, "]")

					if fieldsStart != -1 && fieldsEnd != -1 && fieldsStart < fieldsEnd {
						fieldsContent := indexContent[fieldsStart : fieldsEnd+1]

						var result [][]string

						// Parse the JSON string
						err := json.Unmarshal([]byte(fieldsContent), &result)
						if err != nil {
							log.Fatalf("Error parsing JSON: %v", err)
						}

						currentEntity.CompositeIndex = result

					}

				}
			}
		}

		for _, line := range rawSchema.Fields {
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
						//
						// Check if field is unique or not
						//
						// Extract content between parentheses
						startIdx := strings.Index(line, "(")
						endIdx := strings.Index(line, ")")
						if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
							parenthesesContent := line[startIdx+1 : endIdx]

							// Parse arguments
							if strings.Contains(parenthesesContent, "unique: true") {
								field.IsUnique = true
							}
						}
					} else if strings.Contains(rest, "@unique") {
						field.IsUnique = true
					}
					if strings.Contains(rest, "@derivedFrom") {
						field.DerivedFrom = true
					}
					if !isScalar(fieldType) {
						field.Relation = fieldType
					}
					currentEntity.Fields = append(currentEntity.Fields, field)
				}
			}
		}

		entities = append(entities, *currentEntity)

	}

	return entities, nil

}
