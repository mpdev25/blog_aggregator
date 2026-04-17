-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feed
    FOREIGN KEY(feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

AlTER TABLE posts ALTER COLUMN title DROP NOT NULL;
ALTER TABLE posts DROP CONSTRAINT posts_title_key;
-- +goose Down
DROP TABLE posts;