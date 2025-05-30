-- +goose Up
CREATE TABLE tokens
(
    id         VARCHAR(40) PRIMARY KEY,
    bot_id     VARCHAR(40) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rate       INTEGER   NOT NULL DEFAULT '0',
    period     TEXT,
    FOREIGN KEY (bot_id) REFERENCES bots (id)
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
