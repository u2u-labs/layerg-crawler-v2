package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u2u-labs/layerg-crawler/cmd/query_service"
	"github.com/u2u-labs/layerg-crawler/internal/generator"
	"gopkg.in/yaml.v3"
)

func main() {
	// Define flags for schema input and output directory.
	schemaPath := flag.String("schema", "./schema.graphql", "Path to GraphQL schema file")
	outputDir := flag.String("out", "./generated", "Output directory for generated files")
	queriesDir := flag.String("queries", "./db", "Output directory for queries files")
	flag.Parse()

	// Ensure output directory exists.
	err := os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	err = os.MkdirAll(*queriesDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create queries directory: %v", err)
	}

	// Parse the GraphQL schema.
	entities, enums, err := generator.ParseGraphQLSchema(*schemaPath)
	if err != nil {
		log.Fatalf("failed to parse GraphQL schema: %v", err)
	}

	// Generate Go models.
	if err := generator.GenerateGoModels(entities, *outputDir); err != nil {
		log.Fatalf("failed to generate Go models: %v", err)
	}

	// Generate migration scripts.
	if err := generator.GenerateMigrationScripts(entities, enums, *outputDir); err != nil {
		log.Fatalf("failed to generate migration scripts: %v", err)
	}

	if err := generator.GenerateSQLCQueries(entities, *queriesDir); err != nil {
		log.Fatalf("failed to generate SQLC queries: %v", err)
	}

	if err := query_service.GenerateQueryService(*schemaPath); err != nil {
		log.Fatalf("failed to generate queries for Query service: %v", err)
	}

	// // Double-check the ABI mapping defined in your subgraph YAML.
	// if err := generator.CheckAbiMapping("./abis/erc20.abi.json"); err != nil {
	// 	log.Fatalf("ABI check failed: %v", err)
	// }

	// Load subgraph config
	data, err := os.ReadFile("subgraph.yaml")
	if err != nil {
		log.Fatalf("failed to read subgraph.yaml: %v", err)
	}

	var config generator.CrawlerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("failed to parse subgraph.yaml: %v", err)
	}

	// Generate event handler skeleton.
	if err := generator.GenerateEventHandlers(&config, *outputDir); err != nil {
		log.Fatalf("failed to generate event handlers: %v", err)
	}

	// Generate event router
	if err := generator.GenerateEventRouter(&config, *outputDir); err != nil {
		log.Fatalf("failed to generate event router: %v", err)
	}

	// Generate mappings
	if err := generator.GenerateMappings(&config, *outputDir); err != nil {
		log.Fatalf("failed to generate mappings: %v", err)
	}

	fmt.Println("Code generation completed successfully!")
}
