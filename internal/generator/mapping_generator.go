package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenerateMappings(config *CrawlerConfig, outputDir string) error {
	// Create mappings directory
	mappingsDir := filepath.Join(outputDir, "mappings")
	if err := os.MkdirAll(mappingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create mappings directory: %w", err)
	}

	// Template for mapping interfaces and base structs
	mappingTemplate := `// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// {{ .Name }}Handler defines the interface for handling {{ .Name }} events
type {{ .Name }}Handler interface {
	Handle{{ .Name }}(ctx context.Context, event *eventhandlers.{{ .Name }}) error
}

// Base{{ .Name }}Mapping provides the base structure for {{ .Name }} event mappings
type Base{{ .Name }}Mapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// New{{ .Name }}Mapping creates a new base mapping for {{ .Name }} events
func New{{ .Name }}Mapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *Base{{ .Name }}Mapping {
	return &Base{{ .Name }}Mapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}
`

	// Add title function to template
	funcMap := template.FuncMap{
		"title": strings.Title,
	}

	// Generate mapping files for each event
	events := parseEventsFromConfig(config)
	for _, event := range events {
		fileName := fmt.Sprintf("%s_mapping.go", strings.ToLower(event.Name))
		f, err := os.Create(filepath.Join(mappingsDir, fileName))
		if err != nil {
			return fmt.Errorf("failed to create mapping file %s: %w", fileName, err)
		}

		tmpl, err := template.New("mapping").Funcs(funcMap).Parse(mappingTemplate)
		if err != nil {
			f.Close()
			return fmt.Errorf("failed to parse mapping template: %w", err)
		}

		if err := tmpl.Execute(f, event); err != nil {
			f.Close()
			return fmt.Errorf("failed to execute mapping template: %w", err)
		}
		f.Close()
	}

	return nil
}
