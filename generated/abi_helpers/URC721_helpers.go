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

const URC721 = `[
  {
    "inputs": [],
    "name": "MINT_PRICE",
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
        "internalType": "address",
        "name": "owner",
        "type": "address"
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
    "inputs": [],
    "name": "baseExtension",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
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
        "name": "",
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
    "inputs": [],
    "name": "maxSupply",
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
    "name": "name",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "owner",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
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
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "paused",
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
  },
  {
    "inputs": [],
    "name": "symbol",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
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
    "name": "tokenURI",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "totalSupply",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]`


func URC721MINT_PRICE(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("MINT_PRICE")
	if err != nil {
		sugar.Errorf("Failed to pack data for MINT_PRICE: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "MINT_PRICE", result)
	if err != nil {
		sugar.Errorf("Failed to unpack MINT_PRICE: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC721BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, owner common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
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
	err = contractABI.UnpackIntoInterface(&result, "balanceOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOf: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC721BaseExtension(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (string, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return "", err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("baseExtension")
	if err != nil {
		sugar.Errorf("Failed to pack data for baseExtension: %v", err)
		return "", err
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
		return "", err
	}

	// Unpack the result
	var unpackVar string
	err = contractABI.UnpackIntoInterface(&result, "baseExtension", result)
	if err != nil {
		sugar.Errorf("Failed to unpack baseExtension: %v", err)
		return "", err
	}
	return unpackVar, nil
}


func URC721GetApproved(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
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
	err = contractABI.UnpackIntoInterface(&result, "getApproved", result)
	if err != nil {
		sugar.Errorf("Failed to unpack getApproved: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}


func URC721IsApprovedForAll(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, owner common.Address, operator common.Address) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
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
	err = contractABI.UnpackIntoInterface(&result, "isApprovedForAll", result)
	if err != nil {
		sugar.Errorf("Failed to unpack isApprovedForAll: %v", err)
		return false, err
	}
	return unpackVar, nil
}


func URC721MaxSupply(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("maxSupply")
	if err != nil {
		sugar.Errorf("Failed to pack data for maxSupply: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "maxSupply", result)
	if err != nil {
		sugar.Errorf("Failed to unpack maxSupply: %v", err)
		return nil, err
	}
	return unpackVar, nil
}


func URC721Name(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (string, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return "", err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("name")
	if err != nil {
		sugar.Errorf("Failed to pack data for name: %v", err)
		return "", err
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
		return "", err
	}

	// Unpack the result
	var unpackVar string
	err = contractABI.UnpackIntoInterface(&result, "name", result)
	if err != nil {
		sugar.Errorf("Failed to unpack name: %v", err)
		return "", err
	}
	return unpackVar, nil
}


func URC721Owner(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return common.Address{}, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("owner")
	if err != nil {
		sugar.Errorf("Failed to pack data for owner: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "owner", result)
	if err != nil {
		sugar.Errorf("Failed to unpack owner: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}


func URC721OwnerOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
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
	err = contractABI.UnpackIntoInterface(&result, "ownerOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack ownerOf: %v", err)
		return common.Address{}, err
	}
	return unpackVar, nil
}


func URC721Paused(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return false, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("paused")
	if err != nil {
		sugar.Errorf("Failed to pack data for paused: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "paused", result)
	if err != nil {
		sugar.Errorf("Failed to unpack paused: %v", err)
		return false, err
	}
	return unpackVar, nil
}


func URC721SupportsInterface(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, interfaceId []byte) (bool, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
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
	err = contractABI.UnpackIntoInterface(&result, "supportsInterface", result)
	if err != nil {
		sugar.Errorf("Failed to unpack supportsInterface: %v", err)
		return false, err
	}
	return unpackVar, nil
}


func URC721Symbol(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (string, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return "", err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("symbol")
	if err != nil {
		sugar.Errorf("Failed to pack data for symbol: %v", err)
		return "", err
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
		return "", err
	}

	// Unpack the result
	var unpackVar string
	err = contractABI.UnpackIntoInterface(&result, "symbol", result)
	if err != nil {
		sugar.Errorf("Failed to unpack symbol: %v", err)
		return "", err
	}
	return unpackVar, nil
}


func URC721TokenURI(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (string, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return "", err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("tokenURI", tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for tokenURI: %v", err)
		return "", err
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
		return "", err
	}

	// Unpack the result
	var unpackVar string
	err = contractABI.UnpackIntoInterface(&result, "tokenURI", result)
	if err != nil {
		sugar.Errorf("Failed to unpack tokenURI: %v", err)
		return "", err
	}
	return unpackVar, nil
}


func URC721TotalSupply(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address) (*big.Int, error) {
	// Create ABI instance
	contractABI, err := abi.JSON(strings.NewReader(URC721))
	if err != nil {
		sugar.Errorf("Failed to parse ABI: %v", err)
		return nil, err
	}

	// Prepare the function call data
	data, err := contractABI.Pack("totalSupply")
	if err != nil {
		sugar.Errorf("Failed to pack data for totalSupply: %v", err)
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
	err = contractABI.UnpackIntoInterface(&result, "totalSupply", result)
	if err != nil {
		sugar.Errorf("Failed to unpack totalSupply: %v", err)
		return nil, err
	}
	return unpackVar, nil
}

