package postgres

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/mdma-backend/mdma-backend/internal/api/data"
)

func (db DB) GetAggregatedData(dataType string, meshNodeUUIDs []string, startTime time.Time, endTime time.Time, sampleTime time.Duration, sampleCount int, aggregateFunction string) (data.AggregatedData, error) {
	timeStamps, err := identifyTimeStamps(startTime, endTime, sampleTime, sampleCount)
	if err != nil {
		return data.AggregatedData{}, err
	}

	query, params, err := createQuery(dataType, aggregateFunction)
	if err != nil {
		return data.AggregatedData{}, err
	}

	aggregatedData, err := db.getAggregatedDataSamples(timeStamps, query, params, meshNodeUUIDs)
	if err != nil {
		return data.AggregatedData{}, err
	}

	aggregatedData.AggregateFunction = aggregateFunction
	aggregatedData.DataType = dataType
	aggregatedData.MeshNodeUUIDs = meshNodeUUIDs

	return aggregatedData, nil
}

func identifyTimeStamps(startTime time.Time, endTime time.Time, sampleTime time.Duration, sampleCount int) ([]time.Time, error) {
	if endTime == time.Unix(0, 0) {
		endTime = time.Now()
	}

	duration := endTime.Sub(startTime)
	intervals := sampleCount

	if sampleTime != time.Duration(0) {
		intervals = int(duration / sampleTime)
	} else {
		sampleTime = duration / time.Duration(sampleCount)
	}

	var timeStamps []time.Time
	for i := 0; i <= intervals; i++ {
		timestamp := startTime.Add(sampleTime * time.Duration(i))
		timeStamps = append(timeStamps, timestamp)
	}

	return timeStamps, nil
}

func createQuery(dataType string, aggregateFunction string) (string, []interface{}, error) {
	aggregateFunction = strings.ToLower(aggregateFunction)
	query := "SELECT "

	switch aggregateFunction {
	case "count":
		query += "COUNT(value)"
	case "sum":
		query += "SUM(value::numeric)"
	case "minimum":
		query += "MIN(value)"
	case "maximum":
		query += "MAX(value)"
	case "average":
		query += "AVG(value::numeric)"
	case "range":
		query += "MAX(value::numeric) - MIN(value::numeric)"
	case "median":
		query += "PERCENTILE_DISC(0.5) WITHIN GROUP (ORDER BY d.value) AS value"
	}

	query += `
	FROM data d
	JOIN data_type dt ON d.data_type_id = dt.id
	WHERE dt.name = $1
	`

	params := []interface{}{dataType}

	return query, params, nil
}

func (db DB) getAggregatedDataSamples(timeStamps []time.Time, baseQuery string, baseParams []interface{}, meshNodeUUIDs []string) (data.AggregatedData, error) {
	var aggregatedData data.AggregatedData
	for currentStartTime := 0; currentStartTime < len(timeStamps)-1; currentStartTime++ {

		query := baseQuery
		params := baseParams
		query += " AND d.measured_at >= $" + strconv.Itoa(len(params)+1)
		params = append(params, timeStamps[currentStartTime])

		query += " AND d.measured_at < $" + strconv.Itoa(len(params)+1)
		params = append(params, timeStamps[currentStartTime+1])

		if len(meshNodeUUIDs) > 0 {
			query += ` AND mesh_node_id IN (`
			for controllerNumber, uuid := range meshNodeUUIDs {
				if controllerNumber != 0 {
					query += `, `
				}
				query += `$` + strconv.Itoa(len(params)+1)
				params = append(params, uuid)
			}
			query += `)`
		}

		rows, err := db.pool.Query(query, params...)
		if err != nil {
			println("error in query")
			return data.AggregatedData{}, err
		}
		defer rows.Close()

		foundRows := false

		for rows.Next() {
			var nullableSampleValue sql.NullString
			err := rows.Scan(&nullableSampleValue)
			if err != nil {
				return data.AggregatedData{}, err
			}

			foundRows = true

			var sampleValue string
			if nullableSampleValue.Valid {
				sampleValue = nullableSampleValue.String
			} else {
				sampleValue = "0"
			}

			sample := data.Sample{
				Value:           sampleValue,
				IntervalStartAt: timeStamps[currentStartTime].String(),
				IntervalEndAt:   timeStamps[currentStartTime+1].String(),
			}

			aggregatedData.Samples = append(aggregatedData.Samples, sample)
		}

		if err := rows.Err(); err != nil {
			return data.AggregatedData{}, err
		}

		if !foundRows {
			sample := data.Sample{
				Value:           "0",
				IntervalStartAt: timeStamps[currentStartTime].String(),
				IntervalEndAt:   timeStamps[currentStartTime+1].String(),
			}
			aggregatedData.Samples = append(aggregatedData.Samples, sample)
		}
	}

	return aggregatedData, nil
}

func (db DB) GetManyData(dataType string, meshNodeUUIDs []string, startTime time.Time, endTime time.Time) (data.ManyData, error) {
	var query = `
	SELECT d.id, d.mesh_node_id, d.measured_at, d.value
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
			query += " d.mesh_node_id = $" + strconv.Itoa(len(params)+1)
			params = append(params, uuid)
		}
		query += ")"
	}

	if startTime != time.Unix(0, 0) {
		query += " AND d.measured_at >= $" + strconv.Itoa(len(params)+1)
		params = append(params, startTime)
	}

	if endTime != time.Unix(0, 0) {
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
		var controllerUUID string
		var measurement data.Measurement

		err := rows.Scan(&measurement.UUID, &controllerUUID, &measurement.MeasuredAt, &measurement.Value)
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
	SELECT d.mesh_node_id, dt.name, d.created_at, d.measured_at, d.value 
	FROM data AS d
	JOIN data_type AS dt 
	ON d.data_type_id = dt.id
	WHERE d.id = $1;
	`

	rows, err := db.pool.Query(query, uuid)
	if err != nil {
		return data.Data{}, err
	}
	defer rows.Close()

	var d data.Data
	d.UUID = uuid

	for rows.Next() {
		err := rows.Scan(&d.MeshNodeUUID, &d.Type, &d.CreatedAt, &d.MeasuredAt, &d.Value)
		if err != nil {
			return data.Data{}, err
		}
	}

	if err := rows.Err(); err != nil {
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
