package postgres

import (
	"encoding/json"

	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) RoleByUsername(username string) (types.Role, error) {
	var r types.Role
	var perms []byte
	if err := db.pool.QueryRow(`
SELECT r.id, r.created_at, r.updated_at, r.name, json_agg(rp.permission) AS permissions
FROM role r
JOIN role_permission rp ON r.id = rp.role_id
JOIN user_account ua ON r.id = ua.role_id
WHERE ua.username = $1
GROUP BY r.id;
`, username).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt, &r.Name, &perms); err != nil {
		return r, err
	}

	if err := json.Unmarshal(perms, &r.Permissions); err != nil {
		return r, err
	}

	return r, nil
}

func (db DB) PasswordHashAndSaltByUsername(username string) (types.Hash, types.Salt, error) {
	var hash types.Hash
	var salt types.Salt
	if err := db.pool.QueryRow(`
SELECT password, salt
FROM user_account
WHERE username = $1;
`, username).Scan(&hash, &salt); err != nil {
		return nil, nil, err
	}

	return hash, salt, nil
}
