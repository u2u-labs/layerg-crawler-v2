
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


// Transfer represents the event data for Transfer(address,address,uint256)
type Transfer struct {
	
	From common.Address // address
	
	To common.Address // address
	
	TokenId *big.Int // uint256
	
	Raw *types.Log
}

func UnpackTransfer(log *types.Log) (*Transfer, error) {
	event := new(Transfer)
	event.Raw = log
	var dataOffset int
	
	
	if len(log.Topics) < 2 {
		return nil, fmt.Errorf("missing topic for indexed parameter from")
	}
	
	event.From = common.HexToAddress(log.Topics[1].Hex())
	
	
	
	
	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("missing topic for indexed parameter to")
	}
	
	event.To = common.HexToAddress(log.Topics[2].Hex())
	
	
	
	
	if len(log.Topics) < 4 {
		return nil, fmt.Errorf("missing topic for indexed parameter tokenId")
	}
	
	event.TokenId = new(big.Int).SetBytes(log.Topics[3].Bytes())
	
	
	
	_ = dataOffset
	return event, nil
}

// MetadataUpdate represents the event data for MetadataUpdate(uint256)
type MetadataUpdate struct {
	
	TokenId *big.Int // uint256
	
	Raw *types.Log
}

func UnpackMetadataUpdate(log *types.Log) (*MetadataUpdate, error) {
	event := new(MetadataUpdate)
	event.Raw = log
	var dataOffset int
	
	
		
		if len(log.Data) < dataOffset+32 {
			return nil, fmt.Errorf("insufficient data for non-indexed parameter _tokenId")
		}
		event.TokenId = new(big.Int).SetBytes(log.Data[dataOffset:dataOffset+32])
		dataOffset += 32
		
	
	
	_ = dataOffset
	return event, nil
}


// EventSignatures maps event signatures to their hex representations
var EventSignatures = map[string]string{
	"Transfer(address,address,uint256)": common.HexToHash(KeccakHash("Transfer(address,address,uint256)")).Hex(),
	"MetadataUpdate(uint256)": common.HexToHash(KeccakHash("MetadataUpdate(uint256)")).Hex(),
}

// HandlerRegistry maps event signatures to their handlers
var HandlerRegistry = map[string]EventHandler{
	EventSignatures["Transfer(address,address,uint256)"]: &DefaultHandler{},
	EventSignatures["MetadataUpdate(uint256)"]: &DefaultHandler{},
}

// KeccakHash returns the Keccak256 hash of a string
func KeccakHash(s string) string {
	return common.BytesToHash(crypto.Keccak256([]byte(s))).Hex()
}

// Event signatures
var TransferEventSignature = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex()
var MetadataUpdateEventSignature = crypto.Keccak256Hash([]byte("MetadataUpdate(uint256)")).Hex()
