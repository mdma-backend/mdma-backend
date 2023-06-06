package postgres

import (
	"database/sql"
	"github.com/mdma-backend/mdma-backend/internal/api/service_account"
	"strconv"
	"time"
)

func (db DB) ServiceAccount(id int) (service_account.ServiceAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, name
		FROM service_account
		WHERE id = $1
	`
	rows, err := db.pool.Query(query, id)
	if err != nil {
		return service_account.ServiceAccount{}, err
	}
	defer rows.Close()

	var serviceAccount service_account.ServiceAccount

	for rows.Next() {
		var idString string
		var roleId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var username string

		err := rows.Scan(&idString, &roleId, &createdAt, &updatedAt, &username)
		if err != nil {
			return service_account.ServiceAccount{}, err
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return service_account.ServiceAccount{}, err
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
		serviceAccount = service_account.ServiceAccount{
			ID:        id,
			RoleID:    roleIdInt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			Username:  username,
		}
	}

	if err := rows.Err(); err != nil {
		return service_account.ServiceAccount{}, err
	}

	return serviceAccount, nil

}

func (db DB) AllServiceAccounts() ([]service_account.ServiceAccount, error) {
	query := `
		SELECT id, role_id, created_at, updated_at, name
		FROM service_account
	`
	rows, err := db.pool.Query(query)
	if err != nil {
		return []service_account.ServiceAccount{}, err
	}
	defer rows.Close()

	var result []service_account.ServiceAccount

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
		serviceAccounts := service_account.ServiceAccount{
			ID:        id,
			RoleID:    roleIdInt,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			Username:  username,
		}
		result = append(result, serviceAccounts)
	}

	return result, nil
}

func (db DB) CreateServiceAccount(roleID int, username string) error {
	query := `
		INSERT INTO service_account (role_id, created_at, name,token)
		VALUES ($1, NOW(), $2,$3)
		
	`
	bytes := []byte{97, 98, 99, 100, 101, 102}
	_, err := db.pool.Exec(
		query,
		strconv.Itoa(roleID),
		username,
		bytes,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) UpdateServiceAccount(id int, roleID int, username string) (service_account.ServiceAccount, error) {
	query := `
		UPDATE service_account
		SET role_id = $2, updated_at =$3, name = $4
		WHERE id = $1
	`
	updatedAt := time.Now().Format(time.RFC3339)
	_, err := db.pool.Exec(
		query,
		id,
		roleID,
		updatedAt,
		username,
	)

	if err != nil {
		return service_account.ServiceAccount{}, err
	}

	updatedAccount, err := db.ServiceAccount(id)
	if err != nil {
		return service_account.ServiceAccount{}, err
	}

	return updatedAccount, nil
}
func (db DB) DeleteServiceAccount(id int) error {
	query := `
		DELETE FROM service_account WHERE id = $1;
	`

	_, err := db.pool.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
