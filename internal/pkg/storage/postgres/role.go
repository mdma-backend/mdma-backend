package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) RoleByID(roleID types.RoleID) (types.Role, error) {
	var r types.Role
	var perms []byte
	if err := db.pool.QueryRow(`
SELECT r.id, r.created_at, r.updated_at, r.name, json_agg(rp.permission) AS permissions
FROM role r
JOIN role_permission rp ON r.id = rp.role_id
WHERE r.id = $1
GROUP BY r.id;
`, roleID).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name, &perms); err != nil {
		return r, err
	}

	if err := json.Unmarshal(perms, &r.Permissions); err != nil {
		return r, err
	}

	return r, nil
}

func (db DB) Roles() ([]types.Role, error) {
	rows, err := db.pool.Query(`
SELECT id, created_at, updated_at, name
FROM role;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []types.Role
	for rows.Next() {
		var r types.Role
		if err := rows.Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

func (db DB) CreateRole(role *types.Role) error {
	query := `
WITH new_role AS (
	INSERT INTO role (name)
	VALUES ($1)
	RETURNING id
)
INSERT INTO role_permission (role_id, permission)
SELECT id, unnest(array[
`
	params := []interface{}{role.Name}

	for i, perm := range role.Permissions {
		paramIndex := i + 2
		query += fmt.Sprintf("$%d", paramIndex)
		params = append(params, perm)

		isLastPerm := i == len(role.Permissions)-1
		if isLastPerm {
			break
		}

		query += ", "
	}

	query += `
])::permission
FROM new_role;
`

	_, err := db.pool.Exec(query, params...)
	if err != nil {
		return err
	}

	if err := db.pool.QueryRow(`
SELECT id, created_at
FROM role
WHERE name = $1;
`, role.Name).Scan(&role.ID, &role.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateRole(roleID types.RoleID, role *types.Role) error {
	tx, err := db.pool.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err = tx.QueryRow(`
UPDATE role
SET name = $1,
	updated_at = now()
WHERE id = $2
RETURNING updated_at;
`, role.Name, roleID).Scan(&role.UpdatedAt); err != nil {
		return err
	}

	if _, err = tx.Exec(`
DELETE FROM role_permission
WHERE role_id = $1;
`, roleID); err != nil {
		return err
	}

	query := "INSERT INTO role_permission (role_id, permission) values "
	var params []interface{}

	for i, perm := range role.Permissions {
		paramIndex := len(params) + 1
		query += fmt.Sprintf("($%d, $%d)", paramIndex, paramIndex+1)
		params = append(params, roleID, perm)

		isLastPerm := i == len(role.Permissions)-1
		if isLastPerm {
			break
		}

		query += ", "
	}

	if _, err := tx.Exec(query, params...); err != nil {
		return err
	}

	return tx.Commit()
}

func (db DB) DeleteRole(roleID types.RoleID) error {
	_, err := db.pool.Exec(`
DELETE FROM role
WHERE id = $1;
`, roleID)
	if err != nil {
		return err
	}

	return nil
}
