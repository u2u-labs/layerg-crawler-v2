package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateEventRouter(config *CrawlerConfig, outputDir string) error {
	// Create output directory
	routerDir := filepath.Join(outputDir, "router")
	if err := os.MkdirAll(routerDir, 0755); err != nil {
		return fmt.Errorf("failed to create router directory: %w", err)
	}

	// Generate router code
	routerTemplate := `// Code generated - DO NOT EDIT.
package router

import (
    "context"
    "fmt"

    "github.com/unicornultrafoundation/go-u2u/core/types"
    "go.uber.org/zap"

    "github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
    "github.com/u2u-labs/layerg-crawler/db/graphqldb"
    dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
    "github.com/u2u-labs/layerg-crawler/cmd/handlers"
)

type EventRouter struct {
    handlers map[string]interface{}
    queries  *dbCon.Queries
    gql      *graphqldb.Queries
    logger   *zap.SugaredLogger
    chainID  int32
}

func NewEventRouter(queries *dbCon.Queries, gql *graphqldb.Queries, logger *zap.SugaredLogger, chainID int32) *EventRouter {
    return &EventRouter{
        handlers: make(map[string]interface{}),
        queries:  queries,
        gql:      gql,
        logger:   logger,
        chainID:  chainID,
    }
}

func (r *EventRouter) Route(ctx context.Context, log *types.Log) error {
    if len(log.Topics) == 0 {
        return nil
    }

    signature := log.Topics[0].Hex()
    switch signature {
    {{- range .Events }}
    case eventhandlers.{{ .Name }}EventSignature:
        handler := handlers.New{{ .Name }}Handler(r.queries, r.gql, r.chainID, r.logger)
        event, err := eventhandlers.Unpack{{ .Name }}(log)
        if err != nil {
            return fmt.Errorf("failed to unpack {{ .Name }} event: %w", err)
        }
        return handler.Handle{{ .Name }}(ctx, event)
    {{- end }}
    default:
        r.logger.Debugw("No handler for event signature", "signature", signature)
        return nil
    }
}
`

	// Parse events from config
	events := parseEventsFromConfig(config)

	// Generate router file
	f, err := os.Create(filepath.Join(routerDir, "router.go"))
	if err != nil {
		return fmt.Errorf("failed to create router file: %w", err)
	}
	defer f.Close()

	tmpl, err := template.New("router").Parse(routerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse router template: %w", err)
	}

	if err := tmpl.Execute(f, map[string]interface{}{
		"Events": events,
	}); err != nil {
		return fmt.Errorf("failed to execute router template: %w", err)
	}

	return nil
}
