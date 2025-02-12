// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package graphqldb

import (
	"database/sql"
)

type Collection struct {
	ID      int32          `json:"id"`
	Address string         `json:"address"`
	Type    sql.NullString `json:"type"`
}

type Post struct {
	ID            int32          `json:"id"`
	Title         string         `json:"title"`
	Content       sql.NullString `json:"content"`
	PublishedDate sql.NullTime   `json:"published_date"`
	AuthorID      sql.NullInt32  `json:"author_id"`
}

type Transfer struct {
	ID        int32          `json:"id"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Amount    sql.NullString `json:"amount"`
	Timestamp sql.NullTime   `json:"timestamp"`
}

type User struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Email       sql.NullString `json:"email"`
	CreatedDate sql.NullTime   `json:"created_date"`
	IsActive    sql.NullBool   `json:"is_active"`
	ProfileID   sql.NullInt32  `json:"profile_id"`
}

type UserProfile struct {
	ID        int32          `json:"id"`
	Bio       sql.NullString `json:"bio"`
	AvatarUrl sql.NullString `json:"avatar_url"`
}
