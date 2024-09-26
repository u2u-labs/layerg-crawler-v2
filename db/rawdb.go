package db

import (
	"github.com/unicornultrafoundation/go-u2u/common"
	"gorm.io/gorm"

	"layerg-crawler/types"
)

func mintErc20Asset(db *gorm.DB, contractAddress *common.Address, address *common.Address, balance int) error {
	if err := db.Create(&types.ERC20Asset{
		Owner:        &common.Address{},
		Balance:      "100000",
		BalanceFloat: 100,
	}).Error; err != nil {
		return err
	}
	return nil
}
