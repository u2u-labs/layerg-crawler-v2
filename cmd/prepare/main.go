package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u2u-labs/layerg-crawler/internal/generator"
)

func main() {
	// Define flags for schema input and output directory.
	schemaPath := flag.String("schema", "./schema.graphql", "Path to GraphQL schema file")
	outputDir := flag.String("out", "./generated", "Output directory for generated files")
	flag.Parse()

	// Ensure output directory exists.
	err := os.MkdirAll(*outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// Parse the GraphQL schema.
	entities, err := generator.ParseGraphQLSchema(*schemaPath)
	if err != nil {
		log.Fatalf("failed to parse GraphQL schema: %v", err)
	}

	// Generate Go models.
	if err := generator.GenerateGoModels(entities, *outputDir); err != nil {
		log.Fatalf("failed to generate Go models: %v", err)
	}

	// Generate migration scripts.
	if err := generator.GenerateMigrationScripts(entities, *outputDir); err != nil {
		log.Fatalf("failed to generate migration scripts: %v", err)
	}

	// Double-check the ABI mapping defined in your subgraph YAML.
	if err := generator.CheckAbiMapping("./abis/erc20.abi.json"); err != nil {
		log.Fatalf("ABI check failed: %v", err)
	}

	// Generate event handler skeleton.
	if err := generator.GenerateEventHandlers(*outputDir); err != nil {
		log.Fatalf("failed to generate event handlers: %v", err)
	}

	fmt.Println("Code generation completed successfully!")
}
