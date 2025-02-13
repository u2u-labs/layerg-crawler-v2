package models

import (
	"time"
)

// User represents the User entity.
type User struct {
	Id string `gorm:"primaryKey;not null"`
	Name string `gorm:"index;not null"`
	Email *string `gorm:"index"`
	CreatedDate *time.Time `gorm:""`
	IsActive *bool `gorm:""`
	ProfileID string `gorm:""`
	Profile *UserProfile `gorm:"-"`
	Posts []Post `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

// UserProfile represents the UserProfile entity.
type UserProfile struct {
	Id string `gorm:"primaryKey;not null"`
	Bio *string `gorm:""`
	AvatarUrl *string `gorm:""`
	CreatedAt time.Time `gorm:"not null"`
}

// Post represents the Post entity.
type Post struct {
	Id string `gorm:"primaryKey;not null"`
	Title string `gorm:"not null"`
	Content *string `gorm:""`
	PublishedDate *time.Time `gorm:""`
	AuthorID string `gorm:"index"`
	Author *User `gorm:"-"`
	CreatedAt time.Time `gorm:"not null"`
}

// Collection represents the Collection entity.
type Collection struct {
	Id string `gorm:"primaryKey;not null"`
	Address string `gorm:"index;not null"`
	Type *string `gorm:""`
	CreatedAt time.Time `gorm:"not null"`
}

// Transfer represents the Transfer entity.
type Transfer struct {
	Id string `gorm:"primaryKey;not null"`
	From string `gorm:"not null"`
	To string `gorm:"not null"`
	Amount *string `gorm:""`
	Timestamp *time.Time `gorm:""`
	CreatedAt time.Time `gorm:"not null"`
}

