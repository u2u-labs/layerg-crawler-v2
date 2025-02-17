package generator

import (
	"encoding/json"
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

	// First collect all event names from the handlers in subgraph.yaml
	handlerEvents := make(map[string]bool)
	for _, ds := range config.DataSources {
		for _, h := range ds.Mapping.Handlers {
			if h.Kind == "EthereumHandlerKind.Event" {
				for _, topic := range h.Filter.Topics {
					name, _ := parseEventSignature(topic)
					if name != "" {
						handlerEvents[name] = true
					}
				}
			}
		}
	}

	// Use a map to store unique events by name, but only for events in handlers
	eventMap := make(map[string]Event)
	for _, ds := range config.DataSources {
		for _, abiConfig := range ds.Options.Abis {
			abiPath := abiConfig.File
			if !filepath.IsAbs(abiPath) {
				abiPath = filepath.Join(".", abiPath)
			}

			abiFile, err := os.ReadFile(abiPath)
			if err != nil {
				return fmt.Errorf("failed to read ABI file %s: %w", abiPath, err)
			}

			var abi ABI
			if err := json.Unmarshal(abiFile, &abi); err != nil {
				return fmt.Errorf("failed to parse ABI file %s: %w", abiPath, err)
			}

			// Process events from this ABI
			for _, item := range abi {
				// Only process events that are in the handlers
				if item.Type == "event" && handlerEvents[item.Name] {
					var params []EventParam
					for _, input := range item.Inputs {
						params = append(params, EventParam{
							Name:    input.Name,
							Type:    input.Type,
							Indexed: input.Indexed,
							GoType:  getGoType(input.Type),
						})
					}

					eventMap[item.Name] = Event{
						Name:      item.Name,
						Signature: buildEventSignature(item.Name, item.Inputs),
						Params:    params,
					}
				}
			}
		}
	}

	// Generate mapping files for each event in handlers
	for _, event := range eventMap {
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

// Helper function to parse event signature from topic
// func parseEventSignature(topic string) (string, error) {
// 	parts := strings.Split(topic, "(")
// 	if len(parts) != 2 {
// 		return "", fmt.Errorf("invalid event signature format: %s", topic)
// 	}
// 	return parts[0], nil
// }

// Helper function to build event signature
func buildEventSignature(name string, inputs []ABIInput) string {
	var types []string
	for _, input := range inputs {
		types = append(types, input.Type)
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(types, ","))
}
