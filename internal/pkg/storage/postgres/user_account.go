package postgres

import (
	"database/sql"
	"github.com/mdma-backend/mdma-backend/internal/api/user_account"
	"strconv"
	"time"
)

func (db DB) UserAccount(id int) (user_account.UserAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, username
		FROM user_account
		WHERE id = $1
	`
	rows, err := db.pool.Query(query, id)
	if err != nil {
		return user_account.UserAccount{}, err
	}
	defer rows.Close()

	var userAccount user_account.UserAccount

	for rows.Next() {
		var idString string
		var roleId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var username string

		err := rows.Scan(&idString, &roleId, &createdAt, &updatedAt, &username)
		if err != nil {
			return user_account.UserAccount{}, err
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return user_account.UserAccount{}, err
		}
		var updatedAtString string
		if updatedAt.Valid {
			updatedAtString = updatedAt.String
		}
		var roleIDString string
		if roleId.Valid {
			roleIDString = roleId.String
		}
		roleIdInt, err := strconv.Atoi(roleIDString)
		userAccount = user_account.UserAccount{
			ID:        id,
			RoleID:    roleIdInt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			Username:  username,
			Password:  nil,
		}
	}

	if err := rows.Err(); err != nil {
		return user_account.UserAccount{}, err
	}

	return userAccount, nil

}

func (db DB) AllUserAccounts() ([]user_account.UserAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, username
		FROM user_account
	`
	rows, err := db.pool.Query(query)
	if err != nil {
		return []user_account.UserAccount{}, err
	}
	defer rows.Close()

	var result []user_account.UserAccount

	for rows.Next() {
		var idString string
		var roleId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var username string

		err := rows.Scan(&idString, &roleId, &createdAt, &updatedAt, &username)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		var updatedAtString string
		if updatedAt.Valid {
			updatedAtString = updatedAt.String
		}
		var roleIDString string
		if roleId.Valid {
			roleIDString = roleId.String
		}
		roleIdInt, err := strconv.Atoi(roleIDString)
		accounts := user_account.UserAccount{
			ID:        id,
			RoleID:    roleIdInt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			Username:  username,
			Password:  nil,
		}
		result = append(result, accounts)
	}

	return result, nil
}

func (db DB) CreateUserAccount(roleID int, createdAT string, username string, password []byte) error {
	query := `
		INSERT INTO user_account (role_id, created_at, username, password)
		VALUES ($1, NOW(), $2, $3)
		
	`
	_, err := db.pool.Exec(
		query,
		strconv.Itoa(roleID),
		username,
		password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateUserAccount(id int, roleID int, username string, password []byte) (user_account.UserAccount, error) {
	query := `
		UPDATE user_account
		SET role_id = $2, updated_at =$3, username = $4, password = $5
		WHERE id = $1
	`
	updatedAt := time.Now().Format(time.RFC3339)
	_, err := db.pool.Exec(
		query,
		id,
		roleID,
		updatedAt,
		username,
		password,
	)

	if err != nil {
		return user_account.UserAccount{}, err
	}

	updatedAccount, err := db.UserAccount(id)
	if err != nil {
		return user_account.UserAccount{}, err
	}

	return updatedAccount, nil
}
func (db DB) DeleteUserAccount(id int) error {
	query := `
		DELETE FROM user_account WHERE id = $1;
	`

	_, err := db.pool.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
