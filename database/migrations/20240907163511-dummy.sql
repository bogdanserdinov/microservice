-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE status AS ENUM (
    'pending',
    'success',
    'failed'
);

-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS dummies (
    id            uuid      PRIMARY KEY DEFAULT uuid_generate_v4(),
    status        status    NOT NULL    DEFAULT 'pending',
    description   TEXT      NOT NULL,
    updated_at    TIMESTAMP             DEFAULT NULL,
    created_at    TIMESTAMP NOT NULL    DEFAULT now()
);
-- +migrate Down
DROP TABLE IF EXISTS dummies;
