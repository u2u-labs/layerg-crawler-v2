package generator

import "os"

// GenerateEventHandlers writes an event handler stub file.
func GenerateEventHandlers(outputDir string) error {
	// Create a subdirectory for eventhandlers so that this file is generated in a separate package.
	eventhandlersDir := outputDir + "/eventhandlers"
	if err := os.MkdirAll(eventhandlersDir, os.ModePerm); err != nil {
		return err
	}
	filePath := eventhandlersDir + "/event_handlers.go"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	content := `package eventhandlers

import (
	"fmt"
	"github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/common"
)

// HandleTransaction processes a contract call transaction event.
func HandleTransaction(tx *types.Transaction) {
	// TODO: Parse transaction data and perform event-specific logic.
	fmt.Println("Handling transaction event:", tx.Hash().Hex())
}

// HandleLog processes a blockchain log event.
func HandleLog(log types.Log) {
	// TODO: Parse log data and extract event details.
	fmt.Println("Handling log event. Address:", log.Address.Hex())
}

// EventHandlerMap maps event signatures to their handler functions.
var EventHandlerMap = map[string]interface{}{
	"approve":  HandleTransaction,
	"Transfer": HandleLog,
}
`
	_, err = f.WriteString(content)
	return err
}
