-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS matreshka;

CREATE TYPE matreshka.config_type AS ENUM (
    'plain',
    'verv',
    'minio',
    'pg',
    'nginx',
    'kv'
    );

CREATE TABLE matreshka.configs
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT UNIQUE NOT NULL,
    type       matreshka.config_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE matreshka.configs_content
(
    config_id INTEGER REFERENCES matreshka.configs (id),
    version   TEXT NOT NULL,
    content   TEXT NOT NULL,
    UNIQUE (config_id, version)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE matreshka.configs_content;
DROP TABLE matreshka.configs;
-- +goose StatementEnd
