package postgres

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
