package db

import (
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
