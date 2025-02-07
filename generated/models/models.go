package models

import (
	"time"
)

// User represents the User entity.
type User struct {
	Id string `gorm:"primaryKey;not null"`
	Name string `gorm:"index;not null"`
	Email string `gorm:"index"`
	CreatedDate time.Time `gorm:""`
	IsActive bool `gorm:""`
	Profile *UserProfile `gorm:""`
}

// UserProfile represents the UserProfile entity.
type UserProfile struct {
	Id string `gorm:"primaryKey;not null"`
	Bio string `gorm:""`
	AvatarUrl string `gorm:""`
}

// Post represents the Post entity.
type Post struct {
	Id string `gorm:"primaryKey;not null"`
	Title string `gorm:"not null"`
	Content string `gorm:""`
	PublishedDate time.Time `gorm:""`
	Author *User `gorm:"index"`
}

