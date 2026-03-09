-- +goose Up
-- +goose StatementBegin

CREATE TABLE configs
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT UNIQUE NOT NULL,
    type_name  TEXT        NOT NULL,
    updated_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE configs_values
(
    config_id INTEGER REFERENCES configs (id),
    key       TEXT DEFAULT '' NOT NULL,
    value     TEXT DEFAULT '' NOT NULL,
    version   TEXT DEFAULT '' NOT NULL,
    UNIQUE (config_id, key, version)
);

CREATE TABLE binary_configs
(
    config_id INTEGER REFERENCES configs (id),
    data      BLOB NOT NULL,
    version   TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE configs_values;
DROP TABLE configs;
-- +goose StatementEnd
