package models

import (
	"time"
)

// Item represents the Item entity.
type Item struct {
	Id string `gorm:"primaryKey;uniqueIndex;not null"`
	TokenId string `gorm:"not null"`
	TokenUri string `gorm:"not null"`
	Standard string `gorm:"not null"`
	Balances []Balance `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

// Balance represents the Balance entity.
type Balance struct {
	Id string `gorm:"primaryKey;not null"`
	ItemID string `gorm:"uniqueIndex;not null"`
	Item *Item `gorm:"-"`
	OwnerID string `gorm:"not null"`
	Owner *User `gorm:"-"`
	Value string `gorm:"not null"`
	UpdatedAt string `gorm:"not null"`
	Contract string `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// MetadataUpdateRecord represents the MetadataUpdateRecord entity.
type MetadataUpdateRecord struct {
	Id string `gorm:"primaryKey;not null"`
	TokenId string `gorm:"not null"`
	ActorID string `gorm:"not null"`
	Actor *User `gorm:"-"`
	Timestamp string `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// User represents the User entity.
type User struct {
	Id string `gorm:"primaryKey;not null"`
	Balances []Balance `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

