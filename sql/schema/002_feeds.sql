-- +goose Up
CREATE TABLE feeds (
id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user
    FOREIGN KEY(user_name)
    REFERENCES users(name)
    ON DELETE CASCADE
);

-- +goose DOWN
DROP TABLE feeds;
