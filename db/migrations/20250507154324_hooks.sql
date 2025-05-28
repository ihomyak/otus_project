-- +goose Up
CREATE TABLE hooks
(
    id         VARCHAR(40) PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    rate       INTEGER   NOT NULL DEFAULT '0',
    period     TEXT
);

-- +goose Down
DROP TABLE IF EXISTS hooks;
