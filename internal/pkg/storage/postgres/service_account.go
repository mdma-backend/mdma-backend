package postgres

import (
	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) ServiceAccount(id types.ServiceAccountID) (types.ServiceAccount, error) {
	var sa types.ServiceAccount
	if err := db.pool.QueryRow(`
SELECT id, role_id, created_at, updated_at, name, token
FROM service_account
WHERE id = $1;
`, id).Scan(&sa.ID, &sa.RoleID, &sa.CreatedAt, &sa.UpdatedAt, &sa.Name, &sa.Token); err != nil {
		return sa, err
	}

	return sa, nil
}

func (db DB) AllServiceAccounts() ([]types.ServiceAccount, error) {
	rows, err := db.pool.Query(`
SELECT id, role_id, created_at, updated_at, name
FROM service_account;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serviceAccounts []types.ServiceAccount
	for rows.Next() {
		var sa types.ServiceAccount
		err := rows.Scan(&sa.ID, &sa.RoleID, &sa.CreatedAt, &sa.UpdatedAt, &sa.Name)
		if err != nil {
			return nil, err
		}
		serviceAccounts = append(serviceAccounts, sa)
	}

	return serviceAccounts, nil
}

func (db DB) CreateServiceAccount(sa *types.ServiceAccount) error {
	if err := db.pool.QueryRow(`
INSERT INTO service_account (role_id, name)
VALUES ($1, $2)
RETURNING id, created_at;
`, sa.RoleID, sa.Name).Scan(&sa.ID, &sa.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateServiceAccount(id types.ServiceAccountID, sa *types.ServiceAccount) error {
	if err := db.pool.QueryRow(`
UPDATE service_account
SET role_id = $1, updated_at = now(), name = $2
WHERE id = $3
RETURNING updated_at;
`, sa.RoleID, sa.Name, id).Scan(&sa.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateServiceAccountToken(id types.ServiceAccountID, t types.Token) error {
	if _, err := db.pool.Exec(`
UPDATE service_account
SET updated_at = now(), token = $1
WHERE id = $2;
`, t.Value, id); err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteServiceAccount(id types.ServiceAccountID) error {
	if _, err := db.pool.Exec(`
DELETE FROM service_account
WHERE id = $1;
`, id); err != nil {
		return err
	}

	return nil
}
