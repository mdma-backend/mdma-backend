package postgres

import (
	"fmt"
	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node"
	"time"
)

func (db DB) GetMeshNodes() ([]mesh_node.MeshNode, error) {
	var query = `
		SELECT *
		FROM mesh_nodes AS mn
	`

	rows, err := db.pool.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var result []mesh_node.MeshNode

	for rows.Next() {
		var meshNode mesh_node.MeshNode
		err := rows.Scan(&meshNode.Uuid, &meshNode.Latitude, &meshNode.Longitude, &meshNode.CreatedAt, &meshNode.UpdatedAt, &meshNode.UpdateId)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, meshNode)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (db DB) GetMeshNode(uuid string) (mesh_node.MeshNode, error) {
	query := `
		SELECT * 
		FROM mesh_nodes
		WHERE uuid = $1;
	`

	// Daten aus der Datenbank abrufen
	rows, err := db.pool.Query(query, uuid)
	if err != nil {
		fmt.Println(err)
		return mesh_node.MeshNode{}, err
	}
	defer rows.Close()

	var meshNode mesh_node.MeshNode
	meshNode.Uuid = uuid

	// Fetch the data from the query result
	for rows.Next() {
		err := rows.Scan(&meshNode.Latitude, &meshNode.Longitude, &meshNode.CreatedAt, &meshNode.UpdatedAt, &meshNode.UpdateId)
		if err != nil {
			fmt.Println(err)
			return mesh_node.MeshNode{}, err
		}
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return mesh_node.MeshNode{}, err
	}

	return meshNode, nil
}

func (db DB) PostMeshNode(latitude float32, longitude float32, updateId float32) error {
	query := `
		INSERT INTO mesh_node 
		(latitude, longitude, updateId, createdAt)
		VALUES ($1, $2, $3, $4);
	`

	createdTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", time.Now().String())

	_, err := db.pool.Exec(query, latitude, longitude, updateId, createdTime)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) PostMeshNodeData(data string) error {
	return nil
}

func (db DB) PutMeshNode(uuid string, latitude float32, longitude float32) error {
	query := `
		UPDATE mesh_node 
		SET lat = $2, lng = $3, updated_at = $4
		WHERE uuid = $1;
	`

	updatedTime, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", time.Now().String())

	_, err := db.pool.Exec(query, uuid, latitude, longitude, updatedTime)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) DeleteMeshNode(uuid string) error {
	query := `
		DELETE FROM mesh_node WHERE uuid = $1;
	`

	_, err := db.pool.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}
