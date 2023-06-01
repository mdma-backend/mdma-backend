package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
)

/*
	func (db DB) GetAggregatedData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string, sampleDuration string, sampleCount int, aggregateFunction string) (data.AggregatedData, error) {
		return data.AggregatedData{}, nil
	}
*/

func (db DB) GetManyData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string) (data.ManyData, error) {
	var query = `
			SELECT d.id, d.controller_id, dt.name, d.created_at, d.measured_at, d.value
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
			return data.ManyData{}, err
		}
		query += " AND d.measured_at > $" + strconv.Itoa(len(meshNodeUUIDs)+2)
		params = append(params, startTime)
	}

	if measuredEnd != time.Unix(0, 0).String() {
		endTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", measuredEnd)
		if err != nil {
			return data.ManyData{}, err
		}
		query += " AND d.measured_at < $" + strconv.Itoa(len(meshNodeUUIDs)+3)
		params = append(params, endTime)
	}

	rows, err := db.pool.Query(query, params...)
	if err != nil {
		return data.ManyData{}, err
	}
	defer rows.Close()

	var result data.ManyData
	var currentMeasuredData *data.MeasuredData

	result.DataType = dataType

	for rows.Next() {
		var id string
		var controllerUUID string
		var dataType string
		var createdAt string
		var measuredAt string
		var value string

		err := rows.Scan(&id, &controllerUUID, &dataType, &createdAt, &measuredAt, &value)
		if err != nil {
			return data.ManyData{}, err
		}

		if currentMeasuredData == nil || currentMeasuredData.MeshnodeUUID != controllerUUID {
			if currentMeasuredData != nil {
				result.MeasuredDatas = append(result.MeasuredDatas, *currentMeasuredData)
			}

			currentMeasuredData = &data.MeasuredData{
				MeshnodeUUID: controllerUUID,
			}
		}

		measurement := data.Measurement{
			UUID:       id,
			MeasuredAt: measuredAt,
			Value:      value,
		}

		currentMeasuredData.Measurements = append(currentMeasuredData.Measurements, measurement)
	}

	if currentMeasuredData != nil {
		result.MeasuredDatas = append(result.MeasuredDatas, *currentMeasuredData)
	}

	if err := rows.Err(); err != nil {
		return data.ManyData{}, err
	}
	return result, nil
}

func (db DB) GetData(uuid string) (data.Data, error) {
	query := `
		SELECT d.controller_id, dt.name, d.created_at, d.measured_at, d.value 
		FROM data AS d
		JOIN data_type AS dt 
		ON d.data_type_id = dt.id
		WHERE d.id = $1;
	`

	rows, err := db.pool.Query(query, uuid)
	if err != nil {
		fmt.Println(err)
		return data.Data{}, err
	}
	defer rows.Close()

	var d data.Data
	d.UUID = uuid

	for rows.Next() {
		err := rows.Scan(&d.ControllerUuid, &d.Type, &d.CreatedAt, &d.MeasuredAt, &d.Value)
		if err != nil {
			fmt.Println(err)
			return data.Data{}, err
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return data.Data{}, err
	}

	return d, nil
}

func (db DB) GetTypes() ([]string, error) {
	query := `
		SELECT name FROM data_type;
	`

	rows, err := db.pool.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataTypes := []string{}

	for rows.Next() {
		var dataType string
		err := rows.Scan(&dataType)
		if err != nil {
			return nil, err
		}
		dataTypes = append(dataTypes, dataType)
	}

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
