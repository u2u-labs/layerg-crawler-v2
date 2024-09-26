package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Url  string
	Name string
}

func NewCockroachDbClient(cfg *DbConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Url+"/"+cfg.Name), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//db.AutoMigrate(&types.ERC20Asset{})
	//mintErc20Asset(db, &common.Address{}, &common.Address{}, 100)
	return db, err
}
