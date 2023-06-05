package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node"
	"strconv"
	"strings"
	"time"
)

// GetMeshNodes Funktioniert
func (db DB) GetMeshNodes() ([]mesh_node.MeshNode, error) {
	var query = `
		SELECT *
		FROM mesh_node
	`

	rows, err := db.pool.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var result []mesh_node.MeshNode

	for rows.Next() {

		var id string
		var updateId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var location string

		err := rows.Scan(&id, &updateId, &createdAt, &updatedAt, &location)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if id == "" {
			fmt.Println(err)
			return nil, err
		}

		locationPoint, err := parseLocation(location)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		var updateIdInt int64
		if updateId.Valid {
			updateIdInt, err = strconv.ParseInt(updateId.String, 10, 64)
		}

		var updatedAtString string
		if updatedAt.Valid {
			updatedAtString = updatedAt.String
		}

		meshNode := mesh_node.MeshNode{
			Id:        id,
			Location:  locationPoint,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			UpdateId:  updateIdInt,
		}

		result = append(result, meshNode)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

// GetMeshNode Funktioniert
func (db DB) GetMeshNode(id string) (mesh_node.MeshNode, error) {
	query := `
		SELECT * 
		FROM mesh_node
		WHERE id = $1;
	`

	// Daten aus der Datenbank abrufen
	rows, err := db.pool.Query(query, id)
	if err != nil {
		fmt.Println(err)
		return mesh_node.MeshNode{}, err
	}
	defer rows.Close()

	var meshNode mesh_node.MeshNode

	// Fetch the data from the query result
	for rows.Next() {
		var id string
		var updateId sql.NullString
		var createdAt string
		var updatedAt sql.NullString
		var location string

		err := rows.Scan(&id, &updateId, &createdAt, &updatedAt, &location)
		if err != nil {
			fmt.Println(err)
			return mesh_node.MeshNode{}, err
		}

		if id == "" {
			fmt.Println(err)
			return mesh_node.MeshNode{}, err
		}

		locationPoint, err := parseLocation(location)
		if err != nil {
			fmt.Println(err)
			return mesh_node.MeshNode{}, err
		}

		var updateIdInt int64
		if updateId.Valid {
			updateIdInt, err = strconv.ParseInt(updateId.String, 10, 64)
		}

		var updatedAtString string
		if updatedAt.Valid {
			updatedAtString = updatedAt.String
		}

		meshNode = mesh_node.MeshNode{
			Id:        id,
			Location:  locationPoint,
			CreatedAt: createdAt,
			UpdatedAt: updatedAtString,
			UpdateId:  updateIdInt,
		}
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return mesh_node.MeshNode{}, err
	}

	return meshNode, nil
}

// PostMeshNode Funktioniert
func (db DB) PostMeshNode(meshNode mesh_node.MeshNode) error {
	query := `
		INSERT INTO mesh_node 
		(id, mesh_node_update_id, location)
		VALUES ($1, $2, point($3, $4));
	`

	_, err := db.pool.Exec(query, meshNode.Id, strconv.Itoa(int(meshNode.UpdateId)), meshNode.Location.Lat, meshNode.Location.Lon)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) PostMeshNodeData(meshNodeId string, data data.Data) error {
	query := "SELECT id FROM data_type WHERE name = $1"

	rows, err := db.pool.Query(query, data.Type)

	if err != nil {
		return err
	}
	defer rows.Close()

	var dataTypeId string
	for rows.Next() {

		err := rows.Scan(&dataTypeId)
		if err != nil {
			return err
		}

	}
	if err := rows.Err(); err != nil {
		return err
	}

	query = `
		INSERT INTO data 
		(id, mesh_node_id, data_type_id,measured_at,  value)
		VALUES ($1, $2, $3, $4);
	`

	//createdTime, _ := time.Parse(time.RFC3339, time.Now().String())

	_, err = db.pool.Exec(query, meshNodeId, dataTypeId, data.MeasuredAt, data.Value)
	if err != nil {
		return err
	}

	return nil
}

// PutMeshNode Funktioniert
func (db DB) PutMeshNode(id string, meshNode mesh_node.MeshNode) error {
	query := `
		UPDATE mesh_node 
		SET mesh_node_update_id = $2,  updated_at = $3,  location = point($4, $5)
		WHERE id = $1;
	`

	updatedTime := time.Now().Format(time.RFC3339)

	_, err := db.pool.Exec(query, id, strconv.Itoa(int(meshNode.UpdateId)), updatedTime, meshNode.Location.Lat, meshNode.Location.Lon)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMeshNode Funktioniert
func (db DB) DeleteMeshNode(id string) error {
	query := `
		DELETE FROM mesh_node WHERE id = $1;
	`

	_, err := db.pool.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func parseLocation(location string) (mesh_node.Point, error) {
	location = strings.TrimSpace(location)

	location = strings.Replace(location, "(", "", -1)
	location = strings.Replace(location, ")", "", -1)

	coords := strings.Split(location, ",")

	latFloat, err := strconv.ParseFloat(coords[0], 32)
	if err != nil {
		return mesh_node.Point{}, err
	}

	lonFloat, err := strconv.ParseFloat(coords[1], 32)
	if err != nil {
		return mesh_node.Point{}, err
	}

	point := mesh_node.Point{
		Lat: float32(latFloat),
		Lon: float32(lonFloat),
	}

	return point, nil
}
