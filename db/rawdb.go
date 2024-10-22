package db

import (
	"github.com/u2u-labs/layerg-crawler/types"
	"gorm.io/gorm"

	"github.com/u2u-labs/layerg-crawler/config"
)

func InsertSupportedChains(gdb *gorm.DB) error {
	if err := gdb.Create(config.U2UTestnet).Error; err != nil {
		return err
	}
	if err := gdb.Create(config.U2UMainnet).Error; err != nil {
		return err
	}
	return nil
}

func InsertDefaultContracts(gdb *gorm.DB) error {
	for _, c := range config.DefaultContracts {
		if err := InsertContract(gdb, c); err != nil {
			return err
		}
	}
	return nil
}

func InsertContract(gdb *gorm.DB, contract *types.Contract) error {
	if err := gdb.Create(contract).Error; err != nil {
		return err
	}
	return nil
}
