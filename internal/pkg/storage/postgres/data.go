package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
)

func (db DB) GetAggregatedData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string, sampleDuration string, sampleCount int, aggregateFunction string) (data.AggregatedData, error) {
	aggregateFunction = strings.ToLower(aggregateFunction)
	query := `SELECT `

	if aggregateFunction == "count" {
		query += `COUNT(value)`
	} else if aggregateFunction == "sum" {
		query += `SUM(value::numeric)`
	} else if aggregateFunction == "minimum" {
		query += `MIN(value)`
	} else if aggregateFunction == "maximum" {
		query += `MAX(value)`
	} else if aggregateFunction == "average" {
		query += `AVG(value::numeric)`
	} else if aggregateFunction == "range" {
		query += `MAX(value::numeric) - MIN(value::numeric)`
	}

	params := []interface{}{dataType}

	query += `
	FROM data d
	JOIN data_type dt ON d.data_type_id = dt.id
	WHERE dt.name = $1
	`

	query, err := newFunction(measuredStart, query, params, measuredEnd, meshNodeUUIDs)
	if err != nil {
		return data.AggregatedData{}, err
	}

	rows, err := db.pool.Query(query, params...)
	if err != nil {
		fmt.Println(err)
		return data.AggregatedData{}, err
	}
	defer rows.Close()

	var aggregatedData data.AggregatedData
	aggregatedData.DataType = dataType
	aggregatedData.MeshNodeUUIDs = meshNodeUUIDs
	aggregatedData.AggregateFunction = aggregateFunction

	for rows.Next() {
		var sampleValue string
		err := rows.Scan(&sampleValue)
		if err != nil {
			fmt.Println(err)
			return data.AggregatedData{}, err
		}

		sample := data.Sample{
			Value: sampleValue,
		}

		aggregatedData.Samples = append(aggregatedData.Samples, sample)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return data.AggregatedData{}, err
	}

	return aggregatedData, nil
}

func newFunction(measuredStart string, query string, params []interface{}, measuredEnd string, meshNodeUUIDs []string) (string, error) {
	if measuredStart != time.Unix(0, 0).String() {
		startTime, err := time.Parse(time.RFC3339, measuredStart)
		if err != nil {
			return "", err
		}
		query += " AND d.measured_at > $" + strconv.Itoa(len(params)+1)
		params = append(params, startTime)
	}

	if measuredEnd != time.Unix(0, 0).String() {
		endTime, err := time.Parse(time.RFC3339, measuredEnd)
		if err != nil {
			return "", err
		}
		query += " AND d.measured_at < $" + strconv.Itoa(len(params)+1)
		params = append(params, endTime)
	}

	if len(meshNodeUUIDs) > 0 {
		query += `AND controller_id IN (`
		for i, uuid := range meshNodeUUIDs {
			if i != 0 {
				query += `, `
			}
			query += `$` + strconv.Itoa(len(params)+1+i)
			params = append(params, uuid)
		}
		query += `)`
	}
	return query, nil
}

/*
func (db DB) GetAggregatedData(dataType string, meshNodeUUIDs []string, measuredStart string, measuredEnd string, sampleDuration string, sampleCount int, aggregateFunction string) (data.AggregatedData, error) {
	aggregateFunction = strings.ToLower(aggregateFunction)

	query := `SELECT `
	if aggregateFunction == "range" {
		query += `MAX(value) - MIN(value)`
	} else if aggregateFunction == "count" {
		query += `COUNT(value)`
	} else if aggregateFunction == "minimum" {
		query += `MIN(value)`
	} else if aggregateFunction == "maximum" {
		query += `MAX(value)`
	} else if aggregateFunction == "sum" {
		query += `SUM(value)`
	} else if aggregateFunction == "median" {
		query = `
			percentile_cont(0.5) WITHIN GROUP (ORDER BY value)
		`
	} else if aggregateFunction == "average" {
		query += `AVG(value)`
	}

	params := []interface{}{dataType}

	query += `
	FROM data d
	JOIN data_type dt ON d.data_type_id = dt.id
	WHERE dt.name = $1
	`

	if measuredStart != time.Unix(0, 0).String() {
		startTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", measuredStart)
		if err != nil {
			return data.AggregatedData{}, err
		}
		query += " AND d.measured_at > $2"
		params = append(params, startTime)
	}

	if measuredEnd != time.Unix(0, 0).String() {
		endTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", measuredEnd)
		if err != nil {
			return data.AggregatedData{}, err
		}
		query += " AND d.measured_at < $3" + strconv.Itoa(len(meshNodeUUIDs)+3)
		params = append(params, endTime)
	}

	if len(meshNodeUUIDs) > 0 {
		query += `AND controller_id IN (`
		for i, uuid := range meshNodeUUIDs {
			if i != 0 {
				query += `, `
			}
			query += `$` + strconv.Itoa(i+4)
			params = append(params, uuid)
		}
		query += `)`
	}

	if sampleDuration != "" {
		duration, err := time.ParseDuration(sampleDuration)
		if err != nil {
			fmt.Println(err)
			return data.AggregatedData{}, err
		}

		sampleRange := duration / time.Duration(sampleCount)

		if aggregateFunction != "median" {
			query += ` GROUP BY measured_at / $4`
			params = append(params, sampleRange)
		}
	}

	rows, err := db.pool.Query(query, params...)
	if err != nil {
		fmt.Println(err)
		return data.AggregatedData{}, err
	}
	defer rows.Close()

	var aggregatedData data.AggregatedData
	aggregatedData.DataType = dataType
	aggregatedData.MeshNodeUUIDs = meshNodeUUIDs
	aggregatedData.AggregateFunction = aggregateFunction

	for rows.Next() {
		var sampleValue string
		err := rows.Scan(&sampleValue)
		if err != nil {
			fmt.Println(err)
			return data.AggregatedData{}, err
		}

		sample := data.Sample{
			Value: sampleValue,
		}

		aggregatedData.Samples = append(aggregatedData.Samples, sample)
	}

	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows && aggregateFunction == "median" {
			fmt.Println(err)
			return data.AggregatedData{}, errors.New("median not available")
		}
		fmt.Println(err)
		return data.AggregatedData{}, err
	}

	fmt.Println(err)
	return aggregatedData, nil
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
			query += " d.controller_id = $" + strconv.Itoa(len(params)+1+i)
			params = append(params, uuid)
		}
		query += ")"
	}

	if measuredStart != time.Unix(0, 0).String() {
		startTime, err := time.Parse(time.RFC3339, measuredStart)
		if err != nil {
			fmt.Println(err)
			return data.ManyData{}, err
		}
		query += " AND d.measured_at > $" + strconv.Itoa(len(params)+1)
		params = append(params, startTime)
	}

	if measuredEnd != time.Unix(0, 0).String() {
		endTime, err := time.Parse(time.RFC3339, measuredEnd)
		if err != nil {
			fmt.Println(err)
			return data.ManyData{}, err
		}
		query += " AND d.measured_at < $" + strconv.Itoa(len(params)+1)
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
