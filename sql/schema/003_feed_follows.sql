-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_user
    FOREIGN KEY(user_name)
    REFERENCES users(name)
    ON DELETE CASCADE,
    CONSTRAINT fk_feed
    FOREIGN KEY(feeds_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLEfeed_follows;