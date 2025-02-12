package models

import (
	"time"
)

// User represents the User entity.
type User struct {
	Id int `gorm:"primaryKey;not null"`
	Name string `gorm:"index;not null"`
	Email string `gorm:"index"`
	CreatedDate time.Time `gorm:""`
	IsActive bool `gorm:""`
	ProfileID int `gorm:""`
	Profile *UserProfile `gorm:"-"`
}

// UserProfile represents the UserProfile entity.
type UserProfile struct {
	Id int `gorm:"primaryKey;not null"`
	Bio string `gorm:""`
	AvatarUrl string `gorm:""`
}

// Post represents the Post entity.
type Post struct {
	Id int `gorm:"primaryKey;not null"`
	Title string `gorm:"not null"`
	Content string `gorm:""`
	PublishedDate time.Time `gorm:""`
	AuthorID int `gorm:"index"`
	Author *User `gorm:"-"`
}

// Collection represents the Collection entity.
type Collection struct {
	Id int `gorm:"primaryKey;not null"`
	Address string `gorm:"index;not null"`
	Type string `gorm:""`
}

// Transfer represents the Transfer entity.
type Transfer struct {
	Id int `gorm:"primaryKey;not null"`
	From string `gorm:"not null"`
	To string `gorm:"not null"`
	Amount string `gorm:""`
	Timestamp time.Time `gorm:""`
}

