package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/crypto"
	"go.uber.org/zap"
)

// EventHandler defines the interface for handling blockchain events
type EventHandler interface {
	// HandleEvent processes a log event
	HandleEvent(ctx context.Context, log *types.Log) error
}

// HandlerRegistry maintains a mapping of event signatures to their handlers
type HandlerRegistry struct {
	handlers map[string]EventHandler
	logger   *zap.SugaredLogger
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry(logger *zap.SugaredLogger) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]EventHandler),
		logger:   logger,
	}
}

// RegisterHandler registers an event handler for a specific event signature
func (r *HandlerRegistry) RegisterHandler(eventSignature string, handler EventHandler) {
	// If the signature is already a hash (starts with 0x), use it directly
	var topicHash string
	if strings.HasPrefix(eventSignature, "0x") {
		topicHash = strings.ToLower(eventSignature)
	} else {
		// Convert human readable signature to topic hash
		topicHash = common.BytesToHash(crypto.Keccak256([]byte(eventSignature))).Hex()
	}

	r.handlers[strings.ToLower(topicHash)] = handler

	r.logger.Infow("Registering handler",
		"raw_signature", eventSignature,
		"topic_hash", topicHash,
		"handler", fmt.Sprintf("%T", handler))
}

// RouteEvent routes an event log to its appropriate handler
func (r *HandlerRegistry) RouteEvent(ctx context.Context, log *types.Log) error {
	if len(log.Topics) == 0 {
		return nil
	}

	signature := strings.ToLower(log.Topics[0].Hex())
	handler, exists := r.handlers[signature]

	// Add debug logging
	// r.logger.Debugw("Attempting to route event",
	// 	"incoming_signature", signature,
	// 	"registered_signatures", r.getRegisteredSignatures(),
	// 	"handler_exists", exists,
	// 	"contract", log.Address.Hex(),
	// )

	if !exists {
		return nil
	}

	return handler.HandleEvent(ctx, log)
}

// Helper method to get all registered signatures
func (r *HandlerRegistry) getRegisteredSignatures() []string {
	signatures := make([]string, 0, len(r.handlers))
	for sig := range r.handlers {
		signatures = append(signatures, sig)
	}
	return signatures
}
