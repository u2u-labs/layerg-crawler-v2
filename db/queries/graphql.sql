-- name: CreateUser :one
INSERT INTO "user" (
    id, name, email, createddate, isactive, profile
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;

-- name: GetUserByName :one
SELECT * FROM "user" WHERE name = $1;

-- name: CreateUserProfile :one
INSERT INTO userprofile (
    id, bio, avatarurl
) VALUES ($1, $2, $3)
RETURNING *;

-- name: CreatePost :one
INSERT INTO post (
    id, title, content, publisheddate, author
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPostsByAuthor :many
SELECT * FROM post WHERE author = $1;

-- name: CreateCollection :one
INSERT INTO collection (
    id, address, type
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCollectionByAddress :one
SELECT * FROM collection WHERE address = $1; 