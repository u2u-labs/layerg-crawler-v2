// Code generated by cmd/generate/query-prepare.go; DO NOT EDIT.
package generated

import (
	"github.com/graphql-go/graphql"
	"github.com/u2u-labs/layerg-crawler/cmd/core"
)



var StringFilterType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "StringFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"gte": &graphql.InputObjectFieldConfig{Type: graphql.String},
		"gt":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"eq":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"lt":  &graphql.InputObjectFieldConfig{Type: graphql.String},
		"lte": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})
var BigIntFilterType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "BigIntFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"gte": &graphql.InputObjectFieldConfig{Type: core.BigIntType},
		"gt":  &graphql.InputObjectFieldConfig{Type: core.BigIntType},
		"eq":  &graphql.InputObjectFieldConfig{Type: core.BigIntType},
		"lt":  &graphql.InputObjectFieldConfig{Type: core.BigIntType},
		"lte": &graphql.InputObjectFieldConfig{Type: core.BigIntType},
	},
})
var IDFilterType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "IDFilter",
	Fields: graphql.InputObjectConfigFieldMap{
		"gte": &graphql.InputObjectFieldConfig{Type: graphql.ID},
		"gt":  &graphql.InputObjectFieldConfig{Type: graphql.ID},
		"eq":  &graphql.InputObjectFieldConfig{Type: graphql.ID},
		"lt":  &graphql.InputObjectFieldConfig{Type: graphql.ID},
		"lte": &graphql.InputObjectFieldConfig{Type: graphql.ID},
	},
})

func CreateQueryFields(resolver *core.QueryResolver) graphql.Fields {
	// Pre-declare all types to handle forward references
	
	var BalanceType *graphql.Object
	
	var MetadataUpdateRecordType *graphql.Object
	
	var UserType *graphql.Object
	
	var ItemType *graphql.Object
	

	// Define all input types first since they don't have relationships
	
	var BalanceWhereInputFields = graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{Type: IDFilterType},
		"item": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		"owner": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		"value": &graphql.InputObjectFieldConfig{Type: BigIntFilterType},
		"updatedAt": &graphql.InputObjectFieldConfig{Type: BigIntFilterType},
		"contract": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		
		"AND": &graphql.InputObjectFieldConfig{},
		"OR":  &graphql.InputObjectFieldConfig{},
	}
	var BalanceWhereInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "BalanceWhereInput",
		Fields: BalanceWhereInputFields,
	})
	BalanceWhereInputFields["AND"].Type = graphql.NewList(BalanceWhereInputType)
	BalanceWhereInputFields["OR"].Type = graphql.NewList(BalanceWhereInputType)
	
	var MetadataUpdateRecordWhereInputFields = graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{Type: IDFilterType},
		"tokenId": &graphql.InputObjectFieldConfig{Type: BigIntFilterType},
		"actor": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		"timestamp": &graphql.InputObjectFieldConfig{Type: BigIntFilterType},
		
		"AND": &graphql.InputObjectFieldConfig{},
		"OR":  &graphql.InputObjectFieldConfig{},
	}
	var MetadataUpdateRecordWhereInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "MetadataUpdateRecordWhereInput",
		Fields: MetadataUpdateRecordWhereInputFields,
	})
	MetadataUpdateRecordWhereInputFields["AND"].Type = graphql.NewList(MetadataUpdateRecordWhereInputType)
	MetadataUpdateRecordWhereInputFields["OR"].Type = graphql.NewList(MetadataUpdateRecordWhereInputType)
	
	var UserWhereInputFields = graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{Type: IDFilterType},
		"balances": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		
		"AND": &graphql.InputObjectFieldConfig{},
		"OR":  &graphql.InputObjectFieldConfig{},
	}
	var UserWhereInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "UserWhereInput",
		Fields: UserWhereInputFields,
	})
	UserWhereInputFields["AND"].Type = graphql.NewList(UserWhereInputType)
	UserWhereInputFields["OR"].Type = graphql.NewList(UserWhereInputType)
	
	var ItemWhereInputFields = graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{Type: IDFilterType},
		"tokenId": &graphql.InputObjectFieldConfig{Type: BigIntFilterType},
		"tokenUri": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		"standard": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		"balances": &graphql.InputObjectFieldConfig{Type: StringFilterType},
		
		"AND": &graphql.InputObjectFieldConfig{},
		"OR":  &graphql.InputObjectFieldConfig{},
	}
	var ItemWhereInputType = graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   "ItemWhereInput",
		Fields: ItemWhereInputFields,
	})
	ItemWhereInputFields["AND"].Type = graphql.NewList(ItemWhereInputType)
	ItemWhereInputFields["OR"].Type = graphql.NewList(ItemWhereInputType)
	

	// Now define all object types with their relationships
	
	BalanceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Balance",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				
				"item": &graphql.Field{
					Type: ItemType,
				},
				
				"owner": &graphql.Field{
					Type: UserType,
				},
				
				"value": &graphql.Field{
					Type: core.BigIntType,
				},
				
				"updatedAt": &graphql.Field{
					Type: core.BigIntType,
				},
				
				"contract": &graphql.Field{
					Type: graphql.String,
				},
				
			}
		}),
	})
	
	MetadataUpdateRecordType = graphql.NewObject(graphql.ObjectConfig{
		Name: "MetadataUpdateRecord",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				
				"tokenId": &graphql.Field{
					Type: core.BigIntType,
				},
				
				"actor": &graphql.Field{
					Type: UserType,
				},
				
				"timestamp": &graphql.Field{
					Type: core.BigIntType,
				},
				
			}
		}),
	})
	
	UserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				
				"balances": &graphql.Field{
					Type: graphql.NewList(BalanceType),
				},
				
			}
		}),
	})
	
	ItemType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Item",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				
				"tokenId": &graphql.Field{
					Type: core.BigIntType,
				},
				
				"tokenUri": &graphql.Field{
					Type: graphql.String,
				},
				
				"standard": &graphql.Field{
					Type: graphql.String,
				},
				
				"balances": &graphql.Field{
					Type: graphql.NewList(BalanceType),
				},
				
			}
		}),
	})
	

	return graphql.Fields{
		
		"Balance": &graphql.Field{
			Type: BalanceType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveSingle("Balance", p)
			},
		},
		"Balances": &graphql.Field{
			Type: graphql.NewList(BalanceType),
			Args: graphql.FieldConfigArgument{
				"page":  &graphql.ArgumentConfig{Type: graphql.Int},
				"limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"order": &graphql.ArgumentConfig{Type: graphql.String},
				"where": &graphql.ArgumentConfig{Type: BalanceWhereInputType},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveMultiple("Balance", p)
			},
		},
		
		"MetadataUpdateRecord": &graphql.Field{
			Type: MetadataUpdateRecordType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveSingle("MetadataUpdateRecord", p)
			},
		},
		"MetadataUpdateRecords": &graphql.Field{
			Type: graphql.NewList(MetadataUpdateRecordType),
			Args: graphql.FieldConfigArgument{
				"page":  &graphql.ArgumentConfig{Type: graphql.Int},
				"limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"order": &graphql.ArgumentConfig{Type: graphql.String},
				"where": &graphql.ArgumentConfig{Type: MetadataUpdateRecordWhereInputType},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveMultiple("MetadataUpdateRecord", p)
			},
		},
		
		"User": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveSingle("User", p)
			},
		},
		"Users": &graphql.Field{
			Type: graphql.NewList(UserType),
			Args: graphql.FieldConfigArgument{
				"page":  &graphql.ArgumentConfig{Type: graphql.Int},
				"limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"order": &graphql.ArgumentConfig{Type: graphql.String},
				"where": &graphql.ArgumentConfig{Type: UserWhereInputType},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveMultiple("User", p)
			},
		},
		
		"Item": &graphql.Field{
			Type: ItemType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveSingle("Item", p)
			},
		},
		"Items": &graphql.Field{
			Type: graphql.NewList(ItemType),
			Args: graphql.FieldConfigArgument{
				"page":  &graphql.ArgumentConfig{Type: graphql.Int},
				"limit": &graphql.ArgumentConfig{Type: graphql.Int},
				"order": &graphql.ArgumentConfig{Type: graphql.String},
				"where": &graphql.ArgumentConfig{Type: ItemWhereInputType},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolver.ResolveMultiple("Item", p)
			},
		},
		
	}
}