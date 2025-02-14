package core

import (
	"io/ioutil"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

type Schema struct {
	Types map[string]*ast.ObjectDefinition
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
	s := &Schema{Types: make(map[string]*ast.ObjectDefinition)}
	for _, def := range doc.Definitions {
		if tdef, ok := def.(*ast.ObjectDefinition); ok {
			s.Types[tdef.Name.Value] = tdef
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
