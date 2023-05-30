package postgres

import (
	"fmt"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
)

func (db DB) Data(uuid string) (data.Data, error) {
	// Daten aus der Datenbank abrufen
	rows, err := db.pool.Query(`
		SELECT d.controller_id, dt.name, d.created_at, d.measured_at, d.value 
		FROM data AS d
		JOIN data_type AS dt 
		ON d.data_type_id = dt.id
		WHERE d.id = $1;
	`, uuid)
	if err != nil {
		fmt.Println(err)
		return data.Data{}, err
	}
	defer rows.Close()

	var d data.Data
	d.Uuid = uuid

	// Fetch the data from the query result
	for rows.Next() {
		err := rows.Scan(&d.ControllerUuid, &d.Type, &d.CreatedAt, &d.MeasuredAt, &d.Value)
		if err != nil {
			fmt.Println(err)
			return data.Data{}, err
		}
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return data.Data{}, err
	}

	return d, nil
}

func (db DB) Types() ([]string, error) {
	rows, err := db.pool.Query("SELECT name FROM data_type;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice zum Speichern der Daten
	dataTypes := []string{}

	// Schleife über die Ergebniszeilen
	for rows.Next() {
		var dataType string
		err := rows.Scan(&dataType)
		if err != nil {
			return nil, err
		}
		dataTypes = append(dataTypes, dataType)
	}

	// Fehlerüberprüfung bei Schleifenausführung
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataTypes, nil
}

func (db DB) DeleteData(uuid string) error {
	_, err := db.pool.Exec("DELETE FROM data WHERE id = $1;", uuid)
	if err != nil {
		return err
	}

	return nil
}
