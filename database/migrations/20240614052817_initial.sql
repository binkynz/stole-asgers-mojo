-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE SCHEMA IF NOT EXISTS doc;

CREATE TABLE IF NOT EXISTS doc.campsite (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id integer UNIQUE NOT NULL,
    status text NOT NULL,
    name text NOT NULL,
    region text NOT NULL,
    x float NOT NULL,
    y float NOT NULL
);

CREATE TABLE IF NOT EXISTS doc.hut (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    asset_id integer UNIQUE NOT NULL,
    name text NOT NULL,
    status text NOT NULL,
    region text NOT NULL,
    x float NOT NULL,
    y float NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

DROP SCHEMA IF EXISTS doc CASCADE;

-- +goose StatementEnd
