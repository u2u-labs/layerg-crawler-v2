package generator

import (
	"strings"
)

// CrawlerConfig represents the subgraph.yaml configuration
type CrawlerConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Schema      struct {
		File string `yaml:"file"`
	} `yaml:"schema"`
	Network struct {
		Name     string   `yaml:"name"`
		ChainId  string   `yaml:"chainId"`
		Endpoint []string `yaml:"endpoint"`
	} `yaml:"network"`
	DataSources []struct {
		Kind       string         `yaml:"kind"`
		Options    ContractConfig `yaml:"options"`
		StartBlock int64          `yaml:"startBlock"`
		Mapping    struct {
			Handlers []struct {
				Kind    string `yaml:"kind"`
				Handler string `yaml:"handler"`
				Filter  struct {
					Function string   `yaml:"function,omitempty"`
					Topics   []string `yaml:"topics,omitempty"`
				} `yaml:"filter"`
			} `yaml:"handlers"`
		} `yaml:"mapping"`
	} `yaml:"dataSources"`
	Repository string `yaml:"repository"`
}

type ContractConfig struct {
	Address string      `yaml:"address"`
	Abis    []AbiConfig `yaml:"abis"`
}

type AbiConfig struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

// Event represents a parsed event from the config
type Event struct {
	Name      string       // Event name (e.g., "Transfer")
	Signature string       // Full signature (e.g., "Transfer(address,address,uint256)")
	Params    []EventParam // Event parameters
}

// EventParam represents a parameter in an event
type EventParam struct {
	Name    string
	Type    string
	Indexed bool
	GoType  string
}

// parseEventsFromConfig extracts events from the subgraph config
func parseEventsFromConfig(config *CrawlerConfig) []Event {
	var events []Event
	seenEvents := make(map[string]bool) // To prevent duplicate events

	for _, ds := range config.DataSources {
		for _, h := range ds.Mapping.Handlers {
			// Only process event handlers
			if h.Kind != "EthereumHandlerKind.Event" {
				continue
			}

			for _, topic := range h.Filter.Topics {
				// Skip if we've already processed this event
				if seenEvents[topic] {
					continue
				}
				seenEvents[topic] = true

				// Parse the event signature
				name, params := parseEventSignature(topic)
				if name == "" {
					continue
				}

				events = append(events, Event{
					Name:      name,
					Signature: topic,
					Params:    params,
				})
			}
		}
	}

	return events
}

// parseEventSignature parses an event signature into name and parameters
func parseEventSignature(signature string) (string, []EventParam) {
	// Split the signature into name and parameters
	parts := strings.Split(signature, "(")
	if len(parts) != 2 {
		return "", nil
	}

	name := parts[0]
	paramStr := strings.TrimRight(parts[1], ")")

	// Parse parameters
	var params []EventParam
	if paramStr != "" {
		for _, part := range strings.Split(paramStr, ",") {
			words := strings.Fields(strings.TrimSpace(part))
			if len(words) == 0 {
				continue
			}

			param := EventParam{
				Type:    words[0],
				Indexed: strings.Contains(part, "indexed"),
				GoType:  getGoType(words[0]),
			}

			// Last word is the parameter name (remove 'indexed' if present)
			param.Name = strings.TrimPrefix(words[len(words)-1], "indexed")
			params = append(params, param)
		}
	}

	return name, params
}

// getGoType converts Solidity types to Go types
func getGoType(solidityType string) string {
	switch solidityType {
	case "address":
		return "common.Address"
	case "uint256", "uint128", "uint64":
		return "*big.Int"
	case "uint32":
		return "uint32"
	case "uint16":
		return "uint16"
	case "uint8":
		return "uint8"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "bytes":
		return "[]byte"
	case "bytes32":
		return "[32]byte"
	default:
		if strings.HasPrefix(solidityType, "bytes") {
			return "[]byte"
		}
		return "interface{}"
	}
}
