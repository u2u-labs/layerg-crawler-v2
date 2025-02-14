package core

import (
	"math/big"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var DateType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Date",
	Description: "The custom Date scalar type",
	Serialize: func(value interface{}) interface{} {
		if t, ok := value.(time.Time); ok {
			return t.Format(time.RFC3339)
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		if s, ok := value.(string); ok {
			t, err := time.Parse(time.RFC3339, s)
			if err == nil {
				return t
			}
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		if v, ok := valueAST.(*ast.StringValue); ok {
			t, err := time.Parse(time.RFC3339, v.Value)
			if err == nil {
				return t
			}
		}
		return nil
	},
})

var BigIntType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "BigInt",
	Description: "The `BigInt` scalar type represents non-fractional signed whole numeric values.",
	Serialize: func(value interface{}) interface{} {
		switch v := value.(type) {
		case *big.Int:
			return v.String()
		case string:
			return v
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch v := value.(type) {
		case string:
			n := new(big.Int)
			n, ok := n.SetString(v, 10)
			if ok {
				return n
			}
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			n := new(big.Int)
			n, ok := n.SetString(valueAST.Value, 10)
			if ok {
				return n
			}
		case *ast.IntValue:
			n := new(big.Int)
			n, ok := n.SetString(valueAST.Value, 10)
			if ok {
				return n
			}
		}
		return nil
	},
})
