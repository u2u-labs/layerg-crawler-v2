package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type EventHandlerConfig struct {
	Kind     string   // e.g., "EthereumHandlerKind.Event"
	Handler  string   // e.g., "HandleTransfer"
	Function string   // e.g., "approve"
	Topics   []string // e.g., ["Transfer(address,address,uint256)"]
}

type ABIEvent struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Anonymous bool   `json:"anonymous"`
	Inputs    []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Indexed bool   `json:"indexed"`
	} `json:"inputs"`
}

type ABI []struct {
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Inputs    []ABIInput `json:"inputs"`
	Anonymous bool       `json:"anonymous"`
}

type ABIInput struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed"`
}

const eventHandlerTemplate = `
// Code generated - DO NOT EDIT.
// This file is generated by event_handler_generator.go

package eventhandlers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/crypto"
	"go.uber.org/zap"
)

// EventHandler defines the interface for handling blockchain events
type EventHandler interface {
	HandleEvent(ctx context.Context, log *types.Log, logger *zap.SugaredLogger) error
}

// DefaultHandler is a basic implementation of EventHandler
type DefaultHandler struct{}

func (h *DefaultHandler) HandleEvent(ctx context.Context, log *types.Log, logger *zap.SugaredLogger) error {
	logger.Infow("Default handler called",
		"signature", log.Topics[0].Hex(),
		"contract", log.Address.Hex(),
		"tx", log.TxHash.Hex(),
	)
	return nil
}

{{range .Events}}
// {{.Name}} represents the event data for {{.Signature}}
type {{.Name}} struct {
	{{range .Params}}
	{{title .Name}} {{.GoType}} // {{.Type}}
	{{end}}
	Raw *types.Log
}

func Unpack{{.Name}}(log *types.Log) (*{{.Name}}, error) {
	event := new({{.Name}})
	event.Raw = log
	var dataOffset int
	{{range $index, $param := .Params}}
	{{if .Indexed}}
	if len(log.Topics) < {{add $index 2}} {
		return nil, fmt.Errorf("missing topic for indexed parameter {{.Name}}")
	}
	{{if eq .Type "address"}}
	event.{{title .Name}} = common.HexToAddress(log.Topics[{{add $index 1}}].Hex())
	{{else if eq .Type "uint256"}}
	event.{{title .Name}} = new(big.Int).SetBytes(log.Topics[{{add $index 1}}].Bytes())
	{{end}}
	{{else}}
		{{if eq .Type "string"}}
		event.{{title .Name}} = string(log.Data)
		{{else if eq .Type "uint256[]"}}
		event.{{title .Name}} = []*big.Int{new(big.Int).SetBytes(log.Data)}
		{{else}}
		if len(log.Data) < dataOffset+32 {
			return nil, fmt.Errorf("insufficient data for non-indexed parameter {{.Name}}")
		}
		event.{{title .Name}} = new(big.Int).SetBytes(log.Data[dataOffset:dataOffset+32])
		dataOffset += 32
		{{end}}
	{{end}}
	{{end}}
	_ = dataOffset
	return event, nil
}
{{end}}

// EventSignatures maps event signatures to their hex representations
var EventSignatures = map[string]string{
{{- range .Events}}
	"{{.Signature}}": common.HexToHash(KeccakHash("{{.Signature}}")).Hex(),
{{- end}}
}

// HandlerRegistry maps event signatures to their handlers
var HandlerRegistry = map[string]EventHandler{
{{- range .Events}}
	EventSignatures["{{.Signature}}"]: &DefaultHandler{},
{{- end}}
}

// KeccakHash returns the Keccak256 hash of a string
func KeccakHash(s string) string {
	return common.BytesToHash(crypto.Keccak256([]byte(s))).Hex()
}

// Event signatures
{{- range .Events}}
var {{.Name}}EventSignature = crypto.Keccak256Hash([]byte("{{.Signature}}")).Hex()
{{- end}}
`

func title(s string) string {
	if s == "" {
		return s
	}
	if s[0] == '_' {
		s = s[1:]
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func getEventSignatureFromABI(config *CrawlerConfig, eventName string) (string, error) {
	for _, ds := range config.DataSources {
		for _, abiConfig := range ds.Options.Abis {
			abiPath := abiConfig.File
			if !filepath.IsAbs(abiPath) {
				abiPath = filepath.Join(".", abiPath)
			}
			abiFile, err := os.ReadFile(abiPath)
			if err != nil {
				continue
			}
			var abi ABI
			if err := json.Unmarshal(abiFile, &abi); err != nil {
				continue
			}
			for _, item := range abi {
				if item.Type == "event" && item.Name == eventName {
					var inputs []string
					for _, input := range item.Inputs {
						inputs = append(inputs, input.Type)
					}
					return fmt.Sprintf("%s(%s)", eventName, strings.Join(inputs, ",")), nil
				}
			}
		}
	}
	return "", fmt.Errorf("event %s not found in any ABI", eventName)
}

func GenerateEventHandlers(config *CrawlerConfig, outputDir string) error {
	var handlers []EventHandlerConfig
	eventMap := make(map[string]struct {
		Name      string
		Signature string
		Params    []EventParam
	})

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
			for _, h := range ds.Mapping.Handlers {
				if h.Kind == "EthereumHandlerKind.Event" {
					for _, topic := range h.Filter.Topics {
						name, _ := parseEventSignature(topic)
						if name == "" {
							continue
						}
						for _, item := range abi {
							if item.Type == "event" && item.Name == name {
								var params []EventParam
								for _, input := range item.Inputs {
									params = append(params, EventParam{
										Name:    input.Name,
										Type:    input.Type,
										Indexed: input.Indexed,
										GoType:  getGoType(input.Type),
									})
								}
								var types []string
								for _, input := range item.Inputs {
									types = append(types, input.Type)
								}
								signature := fmt.Sprintf("%s(%s)", name, strings.Join(types, ","))
								eventMap[name] = struct {
									Name      string
									Signature string
									Params    []EventParam
								}{
									Name:      name,
									Signature: signature,
									Params:    params,
								}
								break
							}
						}
					}
				}
				handlers = append(handlers, EventHandlerConfig{
					Kind:     h.Kind,
					Handler:  h.Handler,
					Function: h.Filter.Function,
					Topics:   h.Filter.Topics,
				})
			}
		}
	}

	var events []struct {
		Name      string
		Signature string
		Params    []EventParam
	}
	for _, event := range eventMap {
		events = append(events, event)
	}

	data := struct {
		Handlers []EventHandlerConfig
		Events   []struct {
			Name      string
			Signature string
			Params    []EventParam
		}
		Handler string
	}{
		Handlers: handlers,
		Events:   events,
		Handler:  "DefaultHandler",
	}

	if err := os.MkdirAll(outputDir+"/eventhandlers", 0755); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"title": title,
		"getEventName": func(topics []string) string {
			if len(topics) == 0 {
				return ""
			}
			name, _ := parseEventSignature(topics[0])
			return name
		},
	}

	tmpl, err := template.New("eventHandlers").Funcs(funcMap).Parse(eventHandlerTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(outputDir + "/eventhandlers/event_handlers.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
