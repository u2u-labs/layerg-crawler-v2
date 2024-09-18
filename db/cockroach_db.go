package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DbConfig struct {
	Url string
}

func New(cfg *DbConfig) (*gorm.DB, error) {
	//host := `localhost user=root password= dbname=layerg port=26257 sslmode=disable TimeZone=Asia/Shanghai`
	db, err := gorm.Open(postgres.Open("postgres://root@localhost:26257"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
