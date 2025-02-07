package generator

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

// AbiDefinition represents an ABI entry (simplified).
type AbiDefinition struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Inputs []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Indexed bool   `json:"indexed,omitempty"`
	} `json:"inputs"`
}

// CheckAbiMapping reads an ABI file and ensures it contains the required items.
func CheckAbiMapping(abiFilePath string) error {
	data, err := ioutil.ReadFile(strings.TrimSpace(abiFilePath))
	if err != nil {
		return err
	}

	var abiEntries []AbiDefinition
	if err = json.Unmarshal(data, &abiEntries); err != nil {
		return err
	}

	var foundApprove, foundTransfer bool
	for _, entry := range abiEntries {
		if entry.Type == "function" && entry.Name == "approve" {
			foundApprove = true
		}
		if entry.Type == "event" && entry.Name == "Transfer" {
			foundTransfer = true
		}
	}

	if !foundApprove {
		return errors.New("approve function not found in ABI")
	}
	if !foundTransfer {
		return errors.New("Transfer event not found in ABI")
	}
	return nil
}
