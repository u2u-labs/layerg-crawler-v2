-- +goose Up
-- Migration script generated from GraphQL schema



CREATE TABLE user (
	id VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	email VARCHAR ,
	createddate TIMESTAMP ,
	isactive BOOLEAN ,
	profile VARCHAR 
);


CREATE TABLE userprofile (
	id VARCHAR NOT NULL,
	bio VARCHAR ,
	avatarurl VARCHAR 
);


CREATE TABLE post (
	id VARCHAR NOT NULL,
	title VARCHAR NOT NULL,
	content VARCHAR ,
	publisheddate TIMESTAMP ,
	author VARCHAR 
);


-- +goose Down
