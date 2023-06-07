BEGIN;

-- Drop all tables
/*
DROP TABLE IF EXISTS user_account, service_account, role_permission, role, data, data_type, mesh_node, mesh_node_update CASCADE;
DROP TYPE IF EXISTS permission;

or

DROP SCHEMA IF EXISTS public CASCADE;
*/

CREATE SCHEMA IF NOT EXISTS public;

-- Erstellen der Tabellen

CREATE TABLE role (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) UNIQUE NOT NULL
);

CREATE TYPE permission AS ENUM (
    'mesh_node_create',
    'mesh_node_read',
    'mesh_node_update',
    'mesh_node_delete',

    'mesh_node_update_create',
    'mesh_node_update_read',
    'mesh_node_update_delete',

    'data_create',
    'data_read',
    'data_delete',

    'user_account_create',
    'user_account_read',
    'user_account_update',
    'user_account_delete',

    'service_account_create',
    'service_account_read',
    'service_account_update',
    'service_account_delete',

    'role_create',
    'role_read',
    'role_update',
    'role_delete'
);

CREATE TABLE role_permission (
    role_id BIGINT NOT NULL REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    permission permission NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_account (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    role_id BIGINT REFERENCES role(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    username VARCHAR(120) NOT NULL,
    password BYTEA NOT NULL,
    salt BYTEA NOT NULL
);

CREATE TABLE service_account (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    role_id BIGINT REFERENCES role(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) NOT NULL,
    token BYTEA NOT NULL
);

CREATE TABLE mesh_node_update (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    data BYTEA NOT NULL,
    version VARCHAR(120) NOT NULL
);

CREATE TABLE mesh_node (
    id UUID NOT NULL PRIMARY KEY,
    mesh_node_update_id BIGINT REFERENCES mesh_node_update(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    location POINT NOT NULL
);

CREATE TABLE data_type (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) NOT NULL
);

CREATE TABLE data (
    id UUID NOT NULL PRIMARY KEY,
    mesh_node_id UUID NOT NULL REFERENCES mesh_node(id) ON DELETE CASCADE ON UPDATE CASCADE,
    data_type_id BIGINT NOT NULL REFERENCES data_type(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    measured_at TIMESTAMP,
    value VARCHAR(120) NOT NULL
);

-- Erstellen von Indizes

CREATE INDEX idx_user_account_username ON user_account (username);
CREATE INDEX idx_data_measured_at ON data (measured_at);
CREATE INDEX idx_mesh_node_location ON mesh_node USING GIST (location);
CREATE INDEX idx_mesh_node_update_version ON mesh_node_update (version);

-- Role Examples

INSERT INTO "role" ("name") VALUES ('admin');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'data_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'data_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'data_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'user_account_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'user_account_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'user_account_update');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'user_account_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'service_account_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'service_account_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'service_account_update');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'service_account_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'role_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'role_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'role_update');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'role_delete');

INSERT INTO "user_account" ("role_id", "username", "password", "salt") VALUES ('1', 'H4r4ldD3rH4ck3r', decode('7g5FK95ihY7C7KuZAcyA4r8+pFyg7rFEfadWoWREp1w=', 'base64'), decode('V0P7HrHEEm7Q1F0VLPMbt1BHBFtKLCFJMnhZJivV5Ms=', 'base64'));

COMMIT;