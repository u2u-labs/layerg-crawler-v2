package db

import (
	"github.com/unicornultrafoundation/go-u2u/common"
	"gorm.io/gorm"

	"github.com/u2u-labs/layerg-crawler/config"
	"github.com/u2u-labs/layerg-crawler/types"
)

func mintErc20Asset(gdb *gorm.DB, contractAddress *common.Address, address *common.Address, balance int) error {
	if err := gdb.Create(&types.ERC20Asset{
		Owner:        "0x000000000000000000000000000000000000000000",
		Balance:      "100000",
		BalanceFloat: 100,
	}).Error; err != nil {
		return err
	}
	return nil
}

func InsertSupportedChains(gdb *gorm.DB) error {
	if err := gdb.Create(config.U2UTestnet).Error; err != nil {
		return err
	}
	if err := gdb.Create(config.U2UMainnet).Error; err != nil {
		return err
	}
	return nil
}
