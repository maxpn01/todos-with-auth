package postgres

import (
	_ "github.com/lib/pq"
)

const userSchema = `
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	email BOOLEAN NOT NULL DEFAULT FALSE
	name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	created_at
);`
