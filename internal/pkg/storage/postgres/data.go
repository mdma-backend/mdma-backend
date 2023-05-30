package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
)

func (db DB) GetManyData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string) ([]data.Data, error) {
	var query = `
		SELECT d.controller_id, dt.name, d.created_at, d.measured_at, d.value
		FROM data d
		JOIN data_type dt ON d.data_type_id = dt.id
		WHERE dt.name = $1
	`
	params := []interface{}{dataType}

	if len(meshNodeUUIDs) != 0 {
		query += " AND ("
		for i, uuid := range meshNodeUUIDs {
			if i != 0 {
				query += " OR"
			}
			query += " d.controller_id = $" + strconv.Itoa(i+2)
			params = append(params, uuid)
		}
		query += ")"
	}

	if measuredStart != time.Unix(0, 0).String() {
		startTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", measuredStart)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		query += " AND d.measured_at > $" + strconv.Itoa(len(meshNodeUUIDs)+2)
		params = append(params, startTime)
	}

	if measuredEnd != time.Unix(0, 0).String() {
		endTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", measuredEnd)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		query += " AND d.measured_at < $" + strconv.Itoa(len(meshNodeUUIDs)+3)
		params = append(params, endTime)
	}

	rows, err := db.pool.Query(query, params...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var result []data.Data

	for rows.Next() {
		var d data.Data
		err := rows.Scan(&d.ControllerUuid, &d.Type, &d.CreatedAt, &d.MeasuredAt, &d.Value)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, d)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (db DB) Data(uuid string) (data.Data, error) {
	query := `
		SELECT d.controller_id, dt.name, d.created_at, d.measured_at, d.value 
		FROM data AS d
		JOIN data_type AS dt 
		ON d.data_type_id = dt.id
		WHERE d.id = $1;
	`

	// Daten aus der Datenbank abrufen
	rows, err := db.pool.Query(query, uuid)
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
	query := `
		SELECT name FROM data_type;
	`

	rows, err := db.pool.Query(query)
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
	query := `
		DELETE FROM data WHERE id = $1;
	`

	_, err := db.pool.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}
