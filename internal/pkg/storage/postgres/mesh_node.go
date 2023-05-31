package postgres

import (
	"fmt"
	"strconv"

	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node"
)

func (db DB) GetMeshNodes(dataType string, meshNodeUUIDs []string) ([]mesh_node.MeshNode, error) {

	var query = `
		SELECT *
		FROM mesh_nodes AS mn
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

	rows, err := db.pool.Query(query, params...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var result []mesh_node.MeshNode

	for rows.Next() {
		var d mesh_node.MeshNode
		err := rows.Scan(&d.Type, &d.CreatedAt, &d.Lat, &d.Lng)
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

func (db DB) GetMeshNode(uuid string) (mesh_node.MeshNode, error) {
	query := `
		SELECT * 
		FROM mesh_nodes AS mn
		WHERE mn.id = $1;
	`

	// Daten aus der Datenbank abrufen
	rows, err := db.pool.Query(query, uuid)
	if err != nil {
		fmt.Println(err)
		return mesh_node.MeshNode{}, err
	}
	defer rows.Close()

	var mn mesh_node.MeshNode
	mn.Uuid = uuid

	// Fetch the data from the query result
	for rows.Next() {
		err := rows.Scan(&mn.Type, &mn.CreatedAt, &mn.Lat, &mn.Lng)
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

	return mn, nil
}

func (db DB) PostMeshNode(Type string, Lat float32, Lng float32) error {
	return nil
}
func (db DB) PostMeshNodeData(Type string, Lat float32, Lng float32) error {
	return nil
}

func (db DB) DeleteMeshNode(uuid string) error {
	query := `
		DELETE FROM mesh_node WHERE id = $1;
	`

	_, err := db.pool.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}
