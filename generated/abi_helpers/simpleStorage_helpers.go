package abi_helpers


// Code generated - DO NOT EDIT.
import (
	"context"
	"math/big"
	"strings"

	u2u "github.com/unicornultrafoundation/go-u2u"
	"github.com/unicornultrafoundation/go-u2u/accounts/abi"
	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"go.uber.org/zap"
)

const SIMPLESTORAGE = `[
  {
    "inputs": [],
    "name": "getValue",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "sender",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]`


func SIMPLESTORAGEGetValue(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(SIMPLESTORAGE))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("getValue")
	if err != nil {
		sugar.Errorf("Failed to pack data for getValue: %v", err)
		return nil, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return nil, err
	}

	// Unpack the result
	var unpackVar *big.Int
	err = contractABI.UnpackIntoInterface(&unpackVar, "getValue", result)
	if err != nil {
		sugar.Errorf("Failed to unpack getValue: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func SIMPLESTORAGESender(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(SIMPLESTORAGE))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return common.Address{}, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("sender")
	if err != nil {
		sugar.Errorf("Failed to pack data for sender: %v", err)
		return common.Address{}, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return common.Address{}, err
	}

	// Unpack the result
	var unpackVar common.Address
	err = contractABI.UnpackIntoInterface(&unpackVar, "sender", result)
	if err != nil {
		sugar.Errorf("Failed to unpack sender: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}

