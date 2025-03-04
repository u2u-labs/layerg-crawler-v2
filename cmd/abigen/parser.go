package abigen

import (
	"encoding/json"
	"fmt"
	"os"
)

// ABIFunction represents a single function in the ABI
type ABIFunction struct {
	Inputs          []ABIArgument `json:"inputs"`
	Name            string        `json:"name"`
	Outputs         []ABIArgument `json:"outputs"`
	StateMutability string        `json:"stateMutability"`
	Type            string        `json:"type"`
}

// ABIArgument represents a function argument or return value
type ABIArgument struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

// ParseABIFile parses an ABI JSON file and returns the functions
func ParseABIFile(filename string) ([]ABIFunction, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var functions []ABIFunction
	if err := json.Unmarshal(data, &functions); err != nil {
		return nil, err
	}

	return functions, nil
}

// FilterViewFunctions filters the ABI to only include view functions and returns them as a JSON string
func FilterViewFunctions(functions []ABIFunction) (string, error) {
	viewFunctions := make([]ABIFunction, 0)

	// Filter only view functions
	for _, fn := range functions {
		if fn.Type == "function" && fn.StateMutability == "view" {
			viewFunctions = append(viewFunctions, fn)
		}
	}

	// Convert back to JSON
	jsonBytes, err := json.MarshalIndent(viewFunctions, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal view functions: %v", err)
	}

	return string(jsonBytes), nil
}

// GenerateABIString generates a formatted ABI string for view functions
func GenerateABIString(filename string) (string, error) {
	functions, err := ParseABIFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI file: %v", err)
	}

	abiString, err := FilterViewFunctions(functions)
	if err != nil {
		return "", err
	}

	// Format as a Go string constant
	return fmt.Sprintf("const ABIString = `%s`", abiString), nil
}
