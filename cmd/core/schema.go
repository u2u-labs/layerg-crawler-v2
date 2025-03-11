package core

import (
	"io/ioutil"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type Schema struct {
	Types map[string]*ast.ObjectDefinition
	Enums map[string]*ast.EnumDefinition
}

func LoadSchema(path string) (*Schema, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	src := source.NewSource(&source.Source{Body: b, Name: "schema"})
	doc, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		return nil, err
	}
	s := &Schema{
		Types: make(map[string]*ast.ObjectDefinition),
		Enums: make(map[string]*ast.EnumDefinition),
	}
	for _, def := range doc.Definitions {
		switch tdef := def.(type) {
		case *ast.ObjectDefinition:
			s.Types[tdef.Name.Value] = tdef
		case *ast.EnumDefinition:
			s.Enums[tdef.Name.Value] = tdef
		}
	}
	return s, nil
}

func (s *Schema) GetFieldType(typeName, fieldName string) ast.Type {
	typeObj, exists := s.Types[typeName]
	if !exists {
		return nil
	}

	for _, field := range typeObj.Fields {
		if field.Name.Value == fieldName {
			return field.Type
		}
	}
	return nil
}

func (s *Schema) IsEnum(typeName string) bool {
	_, exists := s.Enums[typeName]
	return exists
}

func (s *Schema) GetEnumValues(typeName string) []string {
	enum, exists := s.Enums[typeName]
	if !exists {
		return nil
	}
	values := make([]string, len(enum.Values))
	for i, v := range enum.Values {
		values[i] = v.Name.Value
	}
	return values
}
