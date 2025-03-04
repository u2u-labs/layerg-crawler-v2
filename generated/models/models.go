package models

import (
	"time"
)

// Value represents the Value entity.
type Value struct {
	Id string `gorm:"primaryKey;not null"`
	Value string `gorm:"not null"`
	Sender string `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

