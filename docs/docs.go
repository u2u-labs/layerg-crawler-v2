// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/asset": {
            "get": {
                "description": "Get all asset collection of the chain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "Get all asset collection of the chain",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Owner Address",
                        "name": "owner",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ERC721",
                            "ERC1155",
                            "ERC20"
                        ],
                        "type": "string",
                        "description": "Asset Type",
                        "name": "asset_type",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/chain": {
            "get": {
                "description": "Get all supported chains",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chains"
                ],
                "summary": "Get all supported chains",
                "responses": {}
            },
            "post": {
                "description": "Add a new chain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chains"
                ],
                "summary": "Add a new chain",
                "parameters": [
                    {
                        "description": "Add a new chain",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.AddChainParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResponseData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/chain/{chain_id}/collection": {
            "get": {
                "description": "Retrieve all asset collections associated with the specified chain ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "Get all asset collections for a specific chain",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chain ID",
                        "name": "chain_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {}
            },
            "post": {
                "description": "Add a new asset collection to the chain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "Add a new asset collection to the chain",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chain ID",
                        "name": "chain_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Asset collection information",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/utils.AddNewAssetParamsSwagger"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/chain/{chain_id}/collection/{collection_address}": {
            "get": {
                "description": "Get all asset collection of the chain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "Get all asset collection of the chain",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chain ID",
                        "name": "chain_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Collection Address",
                        "name": "collection_address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/chain/{chain_id}/collection/{collection_address}/asset": {
            "get": {
                "description": "Get all asset collection of the chain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "asset"
                ],
                "summary": "Get all asset collection of the chain",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Chain ID",
                        "name": "chain_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Collection Address",
                        "name": "collection_address",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Token ID",
                        "name": "token_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Owner Address",
                        "name": "owner",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "db.AddChainParams": {
            "type": "object",
            "properties": {
                "blockTime": {
                    "type": "integer"
                },
                "chain": {
                    "type": "string"
                },
                "chainID": {
                    "type": "integer"
                },
                "explorer": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "latestBlock": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rpcUrl": {
                    "type": "string"
                }
            }
        },
        "db.AssetType": {
            "type": "string",
            "enum": [
                "ERC721",
                "ERC1155",
                "ERC20"
            ],
            "x-enum-varnames": [
                "AssetTypeERC721",
                "AssetTypeERC1155",
                "AssetTypeERC20"
            ]
        },
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "response.ResponseData": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "utils.AddNewAssetParamsSwagger": {
            "type": "object",
            "properties": {
                "chainID": {
                    "type": "integer"
                },
                "collectionAddress": {
                    "type": "string"
                },
                "decimalData": {
                    "type": "integer"
                },
                "initialBlock": {
                    "type": "integer"
                },
                "lastUpdated": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/db.AssetType"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8085",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "Swagger Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
