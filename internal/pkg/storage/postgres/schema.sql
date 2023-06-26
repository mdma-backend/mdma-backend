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
    'role_delete',

    'area_read'
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
    username VARCHAR(120) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    salt BYTEA NOT NULL
);

CREATE TABLE service_account (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    role_id BIGINT REFERENCES role(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(120) UNIQUE NOT NULL,
    token BYTEA
);

CREATE TABLE mesh_node_update (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version VARCHAR(120) UNIQUE NOT NULL,
    data BYTEA NOT NULL
);

CREATE TABLE mesh_node (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    mesh_node_update_id BIGINT REFERENCES mesh_node_update(id) ON DELETE SET NULL ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL
);

CREATE TABLE data_type (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(120) UNIQUE NOT NULL
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
CREATE INDEX idx_mesh_node_update_version ON mesh_node_update (version);

-- Controller 
INSERT INTO "mesh_node" ("id", "mesh_node_update_id", "created_at", "updated_at", "latitude", "longitude") VALUES ('a53b3f71-f073-4578-9557-92fd19d93bb9', NULL, now(), NULL, 1, 1);
INSERT INTO "mesh_node" ("id", "mesh_node_update_id", "created_at", "updated_at", "latitude", "longitude") VALUES ('c33ea7b6-68a7-4bc6-b1e9-0c365db74081', NULL, now(), NULL, 2, 2);
INSERT INTO "mesh_node" ("id", "mesh_node_update_id", "created_at", "updated_at", "latitude", "longitude") VALUES ('f1aef837-04ac-4316-ae1f-0465bc2eb2fa', NULL, now(), NULL, 23, 2);
INSERT INTO "mesh_node" ("id", "mesh_node_update_id", "created_at", "updated_at", "latitude", "longitude") VALUES ('a8957622-acc5-4ddb-bb1f-17e63d3a514f', NULL, now(), NULL, 2, 2);

-- Datatype
INSERT INTO "data_type" ("created_at", "name") VALUES (now(), 'temperature_dummy');
INSERT INTO "data_type" ("created_at", "name") VALUES (now(), 'humidity_dummy');

-- Data Controller 1
INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('c08fd9d6-0ecb-4932-8156-6c31cf885b46', 'a53b3f71-f073-4578-9557-92fd19d93bb9', '1', now(), '2023-06-02T00:46:16+02:00', '12');

INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('3254f1ed-135d-47fc-8acc-4d97862b55a8', 'a53b3f71-f073-4578-9557-92fd19d93bb9', '1', now(), '2023-06-01T00:46:16+02:00', '32');

INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('ee38b09b-692d-4e7c-bed8-287aea55e573', 'a53b3f71-f073-4578-9557-92fd19d93bb9', '1', now(), '2023-05-30T00:46:16+02:00', '34');

INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('4fec4425-5a90-492b-addd-acbedcb6e616', 'a53b3f71-f073-4578-9557-92fd19d93bb9', '2', now(), '2023-05-29T00:46:16+02:00', '12');

-- Data Controller 2
INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('d38f7c02-8477-499e-b2f3-c38bbba0a2dd', 'c33ea7b6-68a7-4bc6-b1e9-0c365db74081', '2', now(), '2023-05-28T00:46:16+02:00', '7');

INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('0f8601a5-e545-4d38-97af-63350e7f99c2', 'c33ea7b6-68a7-4bc6-b1e9-0c365db74081', '1', now(), '2023-05-27T00:46:16+02:00', '23');

INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('9fa74669-1423-4979-9c52-ff34477d263c', 'c33ea7b6-68a7-4bc6-b1e9-0c365db74081', '1', now(), '2023-05-26T00:46:16+02:00', '41');

-- Data Controller 3
INSERT INTO "data" ("id", "mesh_node_id", "data_type_id", "created_at", "measured_at", "value")
VALUES ('d3d5fcff-2eef-4170-9e8e-fb63a5975a42', 'f1aef837-04ac-4316-ae1f-0465bc2eb2fa', '1', now(), '2023-05-25T00:46:16+02:00', '12');

-- Role Examples
INSERT INTO "role" ("name") VALUES ('admin');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_read');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_delete');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update_create');
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'mesh_node_update_read');
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
INSERT INTO "role_permission" ("role_id", "permission") VALUES ('1', 'area_read');

INSERT INTO "role" ("name") VALUES ('förster');

INSERT INTO "user_account" ("role_id", "username", "password", "salt") VALUES ('1', 'H4r4ldD3rH4ck3r', decode('7g5FK95ihY7C7KuZAcyA4r8+pFyg7rFEfadWoWREp1w=', 'base64'), decode('V0P7HrHEEm7Q1F0VLPMbt1BHBFtKLCFJMnhZJivV5Ms=', 'base64'));

COMMIT;