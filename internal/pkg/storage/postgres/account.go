package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mdma-backend/mdma-backend/internal/api/account"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetUser(id int) (account.UserAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, username, password
		FROM user_account
		WHERE id = $1
	`

	var user account.UserAccount

	err := us.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return account.UserAccount{}, fmt.Errorf("User not found")
		}
		return account.UserAccount{}, err
	}

	return user, nil
}

func (us *UserStore) CreateUser(user account.UserAccount) error {
	query := `
		INSERT INTO user_account (role_id, created_at, updated_at, username, password)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := us.db.Exec(
		query,
		user.RoleID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Username,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) UpdateUser(id int, user account.UserAccount) error {
	query := `
		UPDATE user_account
		SET role_id = $1, updated_at = $2, username = $3, password = $4
		WHERE id = $5
	`

	_, err := us.db.Exec(
		query,
		user.RoleID,
		user.UpdatedAt,
		user.Username,
		user.Password,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
func (us *UserStore) DeleteUser(id int) error {
	query := `
		DELETE FROM user_account
		WHERE id = $1
	`

	_, err := us.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
