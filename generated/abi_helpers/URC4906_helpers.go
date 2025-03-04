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

const URC4906 = `[
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "balance",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "tokenId",
        "type": "uint256"
      }
    ],
    "name": "getApproved",
    "outputs": [
      {
        "internalType": "address",
        "name": "operator",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "internalType": "address",
        "name": "operator",
        "type": "address"
      }
    ],
    "name": "isApprovedForAll",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "tokenId",
        "type": "uint256"
      }
    ],
    "name": "ownerOf",
    "outputs": [
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes4",
        "name": "interfaceId",
        "type": "bytes4"
      }
    ],
    "name": "supportsInterface",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]`


func URC4906BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, owner common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC4906))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("balanceOf", owner)
	if err != nil {
		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
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
	err = contractABI.UnpackIntoInterface(&unpackVar, "balanceOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOf: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC4906GetApproved(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC4906))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return common.Address{}, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("getApproved", tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for getApproved: %v", err)
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
	err = contractABI.UnpackIntoInterface(&unpackVar, "getApproved", result)
	if err != nil {
		sugar.Errorf("Failed to unpack getApproved: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}


func URC4906IsApprovedForAll(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, owner common.Address, operator common.Address) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC4906))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return false, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("isApprovedForAll", owner, operator)
	if err != nil {
		sugar.Errorf("Failed to pack data for isApprovedForAll: %v", err)
		return false, err
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
		return false, err
	}

	// Unpack the result
	var unpackVar bool
	err = contractABI.UnpackIntoInterface(&unpackVar, "isApprovedForAll", result)
	if err != nil {
		sugar.Errorf("Failed to unpack isApprovedForAll: %v", err)
		return false, err
	}
	return unpackVar, nil
}


func URC4906OwnerOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC4906))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return common.Address{}, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("ownerOf", tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for ownerOf: %v", err)
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
	err = contractABI.UnpackIntoInterface(&unpackVar, "ownerOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack ownerOf: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}


func URC4906SupportsInterface(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, interfaceId []byte) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC4906))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return false, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("supportsInterface", interfaceId)
	if err != nil {
		sugar.Errorf("Failed to pack data for supportsInterface: %v", err)
		return false, err
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
		return false, err
	}

	// Unpack the result
	var unpackVar bool
	err = contractABI.UnpackIntoInterface(&unpackVar, "supportsInterface", result)
	if err != nil {
		sugar.Errorf("Failed to unpack supportsInterface: %v", err)
		return false, err
	}
	return unpackVar, nil
}

