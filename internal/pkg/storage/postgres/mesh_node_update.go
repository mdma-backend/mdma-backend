package postgres

import "github.com/mdma-backend/mdma-backend/internal/types"

func (db DB) MeshNodeUpdateByID(id types.MeshNodeUpdateID) (types.MeshNodeUpdate, error) {
	var u types.MeshNodeUpdate
	if err := db.pool.QueryRow(`
SELECT id, created_at, version, data
FROM mesh_node_update
WHERE id = $1;
`, id).Scan(&u.ID, &u.CreatedAt, &u.Version, &u.Data); err != nil {
		return u, err
	}

	return u, nil
}

func (db DB) MeshNodeUpdates() ([]types.MeshNodeUpdate, error) {
	rows, err := db.pool.Query(`
SELECT id, created_at, version
FROM mesh_node_update;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []types.MeshNodeUpdate
	for rows.Next() {
		var u types.MeshNodeUpdate
		if err := rows.Scan(&u.ID, &u.CreatedAt, &u.Version); err != nil {
			return nil, err
		}
		updates = append(updates, u)
	}

	return updates, nil
}

func (db DB) CreateMeshNodeUpdate(u *types.MeshNodeUpdate) error {
	if err := db.pool.QueryRow(`
INSERT INTO mesh_node_update (version, data)
VALUES ($1, $2)
RETURNING id, created_at;
`, u.Version, u.Data).Scan(&u.ID, &u.CreatedAt); err != nil {
		return err
	}

	return nil
}
