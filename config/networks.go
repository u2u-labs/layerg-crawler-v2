package config

import (
	"github.com/u2u-labs/layerg-crawler/types"
)

var (
	U2UTestnet = &types.Chain{
		Id:             1,
		Chain:          "U2U",
		Name:           "Nebulas Testnet",
		RpcUrl:         "https://rpc-nebulas-testnet.uniultra.xyz",
		ChainId:        2484,
		TokenContracts: "0xdFAe88F8610a038AFcDF47A5BC77C0963C65087c",
		Explorer:       "https://testnet.u2uscan.xyz/",
		BlockTime:      500,
	}
	U2UMainnet = &types.Chain{
		Id:        2,
		Chain:     "U2U",
		Name:      "Solaris Mainnet",
		RpcUrl:    "https://rpc-mainnet.uniultra.xyz",
		ChainId:   39,
		Explorer:  "https://u2uscan.xyz/",
		BlockTime: 2000,
	}
)
