package postgres

import (
	"database/sql"
	"github.com/mdma-backend/mdma-backend/internal/api/account"
	"strconv"
	"time"
)

func (db DB) UserAccount(id int) (account.UserAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, username
		FROM user_account
		WHERE id = $1
	`
	rows, err := db.pool.Query(query, id)
	if err != nil {
		return account.UserAccount{}, err
	}
	defer rows.Close()

	var userAccount account.UserAccount

	for rows.Next() {
		var idString string
		var roleId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var username string

		err := rows.Scan(&idString, &roleId, &createdAt, &updatedAt, &username)
		if err != nil {
			return account.UserAccount{}, err
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return account.UserAccount{}, err
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
		userAccount = account.UserAccount{
			ID:        id,
			RoleID:    roleIdInt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			Username:  username,
			Password:  nil,
		}
	}

	if err := rows.Err(); err != nil {
		return account.UserAccount{}, err
	}

	return userAccount, nil

}

func (db DB) AllUserAccounts() ([]account.UserAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, username
		FROM user_account
	`
	rows, err := db.pool.Query(query)
	if err != nil {
		return []account.UserAccount{}, err
	}
	defer rows.Close()

	var result []account.UserAccount

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
		accounts := account.UserAccount{
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
	//createdTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", time.Now().String())
	// Execute the SQL statement
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

func (db DB) UpdateUserAccount(id int, roleID int, username string, password []byte) (account.UserAccount, error) {
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
		return account.UserAccount{}, err
	}

	// Fetch the updated user account from the database
	updatedAccount, err := db.UserAccount(id)
	if err != nil {
		return account.UserAccount{}, err
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
