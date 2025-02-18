package models

import (
	"time"
)

// Item represents the Item entity.
type Item struct {
	Id string `gorm:"primaryKey;not null"`
	TokenId string `gorm:"not null"`
	TokenUri string `gorm:"not null"`
	OwnerID string `gorm:"not null"`
	Owner *User `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

// MetadataUpdateRecord represents the MetadataUpdateRecord entity.
type MetadataUpdateRecord struct {
	Id string `gorm:"primaryKey;not null"`
	TokenId string `gorm:"not null"`
	ActorID string `gorm:"not null"`
	Actor *User `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

// User represents the User entity.
type User struct {
	Id string `gorm:"primaryKey;not null"`
	Items []Item `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

