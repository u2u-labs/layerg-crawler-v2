package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"go.uber.org/zap"
)

type Operation struct {
	EntityType string      `json:"entityType"`
	EntityData interface{} `json:"entityData"`
}

type DataProof struct {
	ID             int64       `json:"id"`
	ChainBlockHash string      `json:"chainBlockHash"`
	ProjectID      string      `json:"projectId"`
	Operations     []Operation `json:"operations"`
}

type BaseHandler struct {
	Queries     *db.Queries
	GQL         *graphqldb.Queries
	ChainID     int32
	Logger      *zap.SugaredLogger
	BlockHash   string
	operations  []Operation
	BlockNumber uint64
}

// Helper function to convert snake_case to camelCase
func toCamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

// Convert map keys from snake_case to camelCase
func convertToCamelCase(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range data {
		// Convert nested maps recursively
		if nestedMap, ok := v.(map[string]interface{}); ok {
			v = convertToCamelCase(nestedMap)
		}
		// Convert slice of maps
		if nestedSlice, ok := v.([]interface{}); ok {
			newSlice := make([]interface{}, len(nestedSlice))
			for i, item := range nestedSlice {
				if itemMap, ok := item.(map[string]interface{}); ok {
					newSlice[i] = convertToCamelCase(itemMap)
				} else {
					newSlice[i] = item
				}
			}
			v = newSlice
		}
		result[toCamelCase(k)] = v
	}
	return result
}

func (h *BaseHandler) AddOperation(entityType string, entityData interface{}, blockHash string, blockNumber uint64) {
	// Convert entityData to map first
	jsonBytes, _ := json.Marshal(entityData)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonBytes, &dataMap)

	// Convert to camelCase
	camelCaseData := convertToCamelCase(dataMap)

	h.operations = append(h.operations, Operation{
		EntityType: entityType,
		EntityData: camelCaseData,
	})
	h.BlockHash = blockHash
	h.BlockNumber = blockNumber
}

func (h *BaseHandler) SubmitToDA() error {
	if len(h.operations) == 0 {
		return nil
	}

	daUrl := viper.GetString("DA_URL")
	projectId := viper.GetString("DA_PROJECT_ID")
	proof := DataProof{
		ID:             int64(h.BlockNumber), // You might want to generate this
		ChainBlockHash: h.BlockHash,
		ProjectID:      projectId,
		Operations:     h.operations,
	}

	jsonData, err := json.Marshal(proof)
	if err != nil {
		return fmt.Errorf("failed to marshal proof: %w", err)
	}

	// TODO: Get URL from config
	url := daUrl + "poi/operations/a4b4ceb8-c05e-4a32-96f9-ee19395295a8"
	fmt.Println("data", string(jsonData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit to DA: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("DA submission failed with status: %d", resp.Body)
		return fmt.Errorf("DA submission failed with status: %d", resp.StatusCode)
	}

	// Clear operations after successful submission
	h.operations = nil
	return nil
}
