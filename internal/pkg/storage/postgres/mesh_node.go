package postgres

import (
	"database/sql"
	"errors"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) MeshNodeById(id types.UUID) (types.MeshNode, error) {
	var n types.MeshNode
	if err := db.pool.QueryRow(`
SELECT id, mesh_node_update_id, created_at, updated_at, latitude, longitude
FROM mesh_node
WHERE id = $1;
`, id).Scan(&n.UUID, &n.UpdateID, &n.CreatedAt, &n.UpdatedAt, &n.Latitude, &n.Longitude); err != nil {
		return n, err
	}

	return n, nil
}

func (db DB) MeshNodes() ([]types.MeshNode, error) {
	rows, err := db.pool.Query(`
SELECT id, mesh_node_update_id, created_at, updated_at, latitude, longitude
FROM mesh_node;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meshNodes []types.MeshNode
	for rows.Next() {
		var n types.MeshNode
		if err := rows.Scan(&n.UUID, &n.UpdateID, &n.CreatedAt, &n.UpdatedAt, &n.Latitude, &n.Longitude); err != nil {
			return nil, err
		}
		meshNodes = append(meshNodes, n)
	}

	return meshNodes, nil
}

// PostMeshNode Funktioniert
func (db DB) CreateMeshNode(n *types.MeshNode) error {
	if err := db.pool.QueryRow(`
INSERT INTO mesh_node 
(id, mesh_node_update_id, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING created_at;
`, n.UUID, n.UpdateID, n.Latitude, n.Longitude).Scan(&n.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) CreateMeshNodeData(id types.UUID, data *data.Data) error {
	tx, err := db.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var dataTypeID uint
	if err := db.pool.QueryRow(`
SELECT id
FROM data_type
WHERE name = $1
`, data.Type).Scan(&dataTypeID); errors.Is(err, sql.ErrNoRows) {
		if err := tx.QueryRow(`
INSERT INTO data_type
(name)
VALUES ($1)
RETURNING id;
`, data.Type).Scan(&dataTypeID); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err = tx.QueryRow(`
INSERT INTO data 
(id, mesh_node_id, data_type_id, measured_at, value)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING id, created_at;
`, id, dataTypeID, data.MeasuredAt, data.Value).Scan(&data.UUID, &data.CreatedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (db DB) CreateManyMeshNodeData(id types.UUID, data []data.Data) error {
	tx, err := db.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, d := range data {
		var dataTypeID uint
		if err := db.pool.QueryRow(`
SELECT id
FROM data_type
WHERE name = $1
`, d.Type).Scan(&dataTypeID); errors.Is(err, sql.ErrNoRows) {
			if err := tx.QueryRow(`
INSERT INTO data_type
(name)
VALUES ($1)
RETURNING id;
`, d.Type).Scan(&dataTypeID); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		if err = tx.QueryRow(`
INSERT INTO data 
(id, mesh_node_id, data_type_id, measured_at, value)
VALUES (gen_random_uuid(), $1, $2, $3, $4)
RETURNING id, created_at;
`, id, dataTypeID, d.MeasuredAt, d.Value).Scan(&d.UUID, &d.CreatedAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db DB) UpdateMeshNode(id types.UUID, n *types.MeshNode) error {
	if err := db.pool.QueryRow(`
UPDATE mesh_node 
SET mesh_node_update_id = $1,  updated_at = now(),  latitude = $2, longitude = $3
WHERE id = $4
RETURNING created_at, updated_at;
`, n.UpdateID, n.Latitude, n.Longitude, id).Scan(&n.CreatedAt, &n.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteMeshNode(id types.UUID) error {
	if _, err := db.pool.Exec(`
DELETE FROM mesh_node WHERE id = $1;
`, id); err != nil {
		return err
	}

	return nil
}
