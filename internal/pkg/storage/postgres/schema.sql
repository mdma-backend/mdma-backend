BEGIN;

-- Drop all tables
/*
DROP TABLE IF EXISTS user_account, service_account, role_permission, role, data, data_type, controller, update CASCADE;
DROP TYPE IF EXISTS permission;

or

DROP SCHEMA IF EXISTS public CASCADE;
*/

CREATE SCHEMA IF NOT EXISTS public;

-- Erstellen der Tabellen

CREATE TABLE role (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) NOT NULL
);

CREATE TYPE permission AS ENUM (
    'controller_create',
    'controller_read',
    'controller_update',
    'controller_delete',
    'sensor_data_create',
    'sensor_data_read',
    'sensor_data_delete',
    'user_account_create',
    'user_account_read',
    'user_account_update',
    'user_account_delete',
    'controller_update_create',
    'controller_update_read',
    'controller_update_update',
    'controller_update_delete',
    'metric_create',
    'metric_read'
);

CREATE TABLE role_permission (
    role_id BIGINT NOT NULL REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    permission permission NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_account (
    id SERIAL NOT NULL PRIMARY KEY,
    role_id BIGINT REFERENCES role(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    username VARCHAR(120) NOT NULL,
    password BYTEA NOT NULL
);

CREATE TABLE service_account (
    id SERIAL NOT NULL PRIMARY KEY,
    role_id BIGINT REFERENCES role(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) NOT NULL,
    token BYTEA NOT NULL
);

CREATE TABLE update (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    data BYTEA NOT NULL,
    version VARCHAR(120) NOT NULL
);

CREATE TABLE controller (
    id UUID NOT NULL PRIMARY KEY,
    update_id BIGINT REFERENCES update(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    location POINT NOT NULL
);

CREATE TABLE data_type (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) NOT NULL
);

CREATE TABLE data (
    id UUID NOT NULL PRIMARY KEY,
    controller_id UUID NOT NULL REFERENCES controller(id) ON DELETE CASCADE ON UPDATE CASCADE,
    data_type_id BIGINT NOT NULL REFERENCES data_type(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    measured_at TIMESTAMP,
    value VARCHAR(120) NOT NULL
);

-- Erstellen von Indizes

CREATE INDEX idx_user_account_username ON user_account (username);
CREATE INDEX idx_data_measured_at ON data (measured_at);
CREATE INDEX idx_controller_location ON controller USING GIST (location);
CREATE INDEX idx_update_version ON update (version);

-- Controller 
INSERT INTO "controller" ("id", "update_id", "created_at", "updated_at", "location") VALUES ('a53b3f71-f073-4578-9557-92fd19d93bb9', NULL, now(), NULL, '1,1');
INSERT INTO "controller" ("id", "update_id", "created_at", "updated_at", "location") VALUES ('c33ea7b6-68a7-4bc6-b1e9-0c365db74081', NULL, now(), NULL, '2,2');

-- Datatype
INSERT INTO "data_type" ("created_at", "updated_at", "name") VALUES (now(), NULL, 'temperature_dummy');
INSERT INTO "data_type" ("created_at", "updated_at", "name") VALUES (now(), NULL, 'humidity_dummy');

-- Data
INSERT INTO "data" ("id", "controller_id", "data_type_id", "created_at", "measured_at", "value") 
VALUES ('c08fd9d6-0ecb-4932-8156-6c31cf885b46', 'a53b3f71-f073-4578-9557-92fd19d93bb9', '1', now(), NULL, '12');

INSERT INTO "data" ("id", "controller_id", "data_type_id", "created_at", "measured_at", "value") 
VALUES ('d38f7c02-8477-499e-b2f3-c38bbba0a2dd', 'c33ea7b6-68a7-4bc6-b1e9-0c365db74081', '2', now(), NULL, 'very wet');

COMMIT;