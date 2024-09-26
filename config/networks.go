package config

import (
	"layerg-crawler/types"
)

var (
	U2UTestnet = &types.Network{
		Chain:    "U2U",
		Name:     "Nebulas Testnet",
		RpcUrl:   "https://rpc-nebulas-testnet.uniultra.xyz/",
		ChainId:  2484,
		Explorer: "https://testnet.u2uscan.xyz/",
	}
	U2UMainnet = &types.Network{
		Chain:    "U2U",
		Name:     "Solaris Mainnet",
		RpcUrl:   "https://rpc-mainnet.uniultra.xyz",
		ChainId:  39,
		Explorer: "https://u2uscan.xyz/",
	}
)
