// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package graphqldb

import (
	"time"
)

type Balance struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	OwnerID   string    `json:"owner_id"`
	Value     string    `json:"value"`
	UpdatedAt string    `json:"updated_at"`
	Contract  string    `json:"contract"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID        string    `json:"id"`
	TokenID   string    `json:"token_id"`
	TokenUri  string    `json:"token_uri"`
	Standard  string    `json:"standard"`
	CreatedAt time.Time `json:"created_at"`
}

type MetadataUpdateRecord struct {
	ID        string    `json:"id"`
	TokenID   string    `json:"token_id"`
	ActorID   string    `json:"actor_id"`
	Timestamp string    `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
