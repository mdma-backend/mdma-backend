package postgres

import (
	"database/sql"
	"errors"

	"github.com/mdma-backend/mdma-backend/internal/types"
)

func (db DB) UserAccountByUsername(username string) (types.UserAccount, error) {
	var ua types.UserAccount
	if err := db.pool.QueryRow(`
SELECT id, role_id, created_at, updated_at, username
FROM user_account
WHERE username = $1;
`, username).Scan(&ua.ID, &ua.RoleID, &ua.CreatedAt, &ua.UpdatedAt, &ua.Username); errors.Is(err, sql.ErrNoRows) {
		return ua, types.ErrNotFound
	} else if err != nil {
		return ua, err
	}

	return ua, nil
}

func (db DB) UserAccountByID(id types.UserAccountID) (types.UserAccount, error) {
	var ua types.UserAccount
	if err := db.pool.QueryRow(`
SELECT id, role_id, created_at, updated_at, username
FROM user_account
WHERE id = $1;
`, id).Scan(&ua.ID, &ua.RoleID, &ua.CreatedAt, &ua.UpdatedAt, &ua.Username); errors.Is(err, sql.ErrNoRows) {
		return ua, types.ErrNotFound
	} else if err != nil {
		return ua, err
	}

	return ua, nil
}

func (db DB) AllUserAccounts() ([]types.UserAccount, error) {
	rows, err := db.pool.Query(`
SELECT id, role_id, created_at, updated_at, username
FROM user_account;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAccounts []types.UserAccount
	for rows.Next() {
		var ua types.UserAccount
		err := rows.Scan(&ua.ID, &ua.RoleID, &ua.CreatedAt, &ua.UpdatedAt, &ua.Username)
		if err != nil {
			return nil, err
		}
		userAccounts = append(userAccounts, ua)
	}

	return userAccounts, nil
}

func (db DB) CreateUserAccount(ua *types.UserAccount, h types.Hash, s types.Salt) error {
	if err := db.pool.QueryRow(`
INSERT INTO user_account (role_id, username, password, salt)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;
`, ua.RoleID, ua.Username, h, s).Scan(&ua.ID, &ua.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateUserAccountPassword(id types.UserAccountID, h types.Hash, s types.Salt) error {
	res, err := db.pool.Exec(`
UPDATE user_account
SET updated_at = now(), password = $1, salt = $2
WHERE id = $3
RETURNING updated_at;
`, h, s, id)
	if err != nil {
		return err
	}

	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return types.ErrNotFound
	}

	return nil
}

func (db DB) UpdateUserAccount(id types.UserAccountID, ua *types.UserAccount) error {
	if err := db.pool.QueryRow(`
UPDATE user_account
SET role_id = $1, updated_at = now(), username = $2
WHERE id = $3
RETURNING updated_at;
`, ua.RoleID, ua.Username, id).Scan(&ua.UpdatedAt); errors.Is(err, sql.ErrNoRows) {
		return types.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteUserAccount(id types.UserAccountID) error {
	res, err := db.pool.Exec(`
DELETE FROM user_account
WHERE id = $1;
`, id)
	if err != nil {
		return err
	}

	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return types.ErrNotFound
	}

	return nil
}
