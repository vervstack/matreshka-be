-- +goose Up
-- +goose StatementBegin
CREATE TYPE config_type AS ENUM (
    'plain',
    'verv',
    'minio',
    'pg',
    'nginx',
    'kv'
    );

CREATE TABLE configs
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT UNIQUE NOT NULL,
    type       config_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE configs_content
(
    config_id INTEGER REFERENCES configs (id),
    version   TEXT NOT NULL,
    content   TEXT NOT NULL,
    UNIQUE (config_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE configs_content;
DROP TABLE configs;
-- +goose StatementEnd
