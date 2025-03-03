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

const URC1155 = `[
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "account",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "id",
        "type": "uint256"
      }
    ],
    "name": "balanceOf",
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
    "inputs": [
      {
        "internalType": "address[]",
        "name": "accounts",
        "type": "address[]"
      },
      {
        "internalType": "uint256[]",
        "name": "ids",
        "type": "uint256[]"
      }
    ],
    "name": "balanceOfBatch",
    "outputs": [
      {
        "internalType": "uint256[]",
        "name": "",
        "type": "uint256[]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "account",
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
  }
]`


func URC1155BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, account common.Address, id *big.Int) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC1155))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("balanceOf", account, id)
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
	err = contractABI.UnpackIntoInterface(&result, "balanceOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOf: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC1155BalanceOfBatch(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, accounts interface{}, ids *big.Int) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC1155))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("balanceOfBatch", accounts, ids)
	if err != nil {
		sugar.Errorf("Failed to pack data for balanceOfBatch: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "balanceOfBatch", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOfBatch: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC1155IsApprovedForAll(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, account common.Address, operator common.Address) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC1155))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return false, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("isApprovedForAll", account, operator)
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
	err = contractABI.UnpackIntoInterface(&result, "isApprovedForAll", result)
	if err != nil {
		sugar.Errorf("Failed to unpack isApprovedForAll: %v", err)
		return false, err
	}
	return unpackVar, nil
}

