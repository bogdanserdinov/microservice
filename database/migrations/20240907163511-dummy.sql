-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE
    OR REPLACE FUNCTION dummys_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TYPE tx_status AS ENUM (
    'pending',
    'success',
    'failed'
);

-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS dummys (
    id            uuid      PRIMARY KEY DEFAULT uuid_generate_v4(),
    status        tx_status NOT NULL    DEFAULT 'pending',
    description   TEXT      NOT NULL,
    updated_at    TIMESTAMP             DEFAULT NULL,
    created_at    TIMESTAMP NOT NULL    DEFAULT now()
);

CREATE TRIGGER update_dummys_modtime BEFORE
    UPDATE ON dummys FOR EACH ROW EXECUTE PROCEDURE dummys_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS dummys_update_updated_at_column ON dummys;
DROP TABLE IF EXISTS dummys;
DROP FUNCTION IF EXISTS dummys_update_updated_at_column();
