-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE
    OR REPLACE FUNCTION dummies_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TYPE status AS ENUM (
    'pending',
    'success',
    'failed'
);

-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS dummies (
    id            uuid      PRIMARY KEY DEFAULT uuid_generate_v4(),
    status        status NOT NULL    DEFAULT 'pending',
    description   TEXT      NOT NULL,
    updated_at    TIMESTAMP             DEFAULT NULL,
    created_at    TIMESTAMP NOT NULL    DEFAULT now()
);

CREATE TRIGGER update_dummies_modtime BEFORE
    UPDATE ON dummies FOR EACH ROW EXECUTE PROCEDURE dummys_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS dummies_update_updated_at_column ON dummies;
DROP TABLE IF EXISTS dummies;
DROP FUNCTION IF EXISTS dummies_update_updated_at_column();
