package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) RoleByUserAccountID(uaId types.UserAccountID) (types.Role, error) {
	var r types.Role
	var perms []byte
	if err := db.pool.QueryRow(`
SELECT r.id, r.created_at, r.updated_at, r.name, json_agg(rp.permission) AS permissions
FROM role r
JOIN role_permission rp ON r.id = rp.role_id
JOIN user_account ua ON r.id = ua.role_id
WHERE ua.id = $1
GROUP BY r.id;
`, uaId).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name, &perms); err != nil {
		return r, err
	}

	if err := json.Unmarshal(perms, &r.Permissions); err != nil {
		return r, err
	}

	return r, nil
}

func (db DB) RoleByServiceAccountID(saId types.ServiceAccountID) (types.Role, error) {
	var r types.Role
	var perms []byte
	if err := db.pool.QueryRow(`
SELECT r.id, r.created_at, r.updated_at, r.name, json_agg(rp.permission) AS permissions
FROM role r
JOIN role_permission rp ON r.id = rp.role_id
JOIN service_account sa ON r.id = sa.role_id
WHERE sa.id = $1
GROUP BY r.id;
`, saId).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name, &perms); err != nil {
		return r, err
	}

	if err := json.Unmarshal(perms, &r.Permissions); err != nil {
		return r, err
	}

	return r, nil
}

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
	tx, err := db.pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(`
INSERT INTO role (name)
VALUES ($1)
RETURNING id, created_at;
`, role.Name).Scan(&role.ID, &role.CreatedAt); err != nil {
		return err
	}

	query := "INSERT INTO role_permission (role_id, permission) values "
	var params []interface{}

	for i, perm := range role.Permissions {
		paramIndex := len(params) + 1
		query += fmt.Sprintf("($%d, $%d)", paramIndex, paramIndex+1)
		params = append(params, role.ID, perm)

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

func (db DB) UpdateRole(roleID types.RoleID, role *types.Role) error {
	tx, err := db.pool.Begin()
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
