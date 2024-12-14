package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/api/internal/logger"
)

func MarshalToJSON(value interface{}) (string, error) {
	linksJSON, err := json.Marshal(value)
	if err != nil {
		logger.Log.Errorf("Failed to marshal value: %v", err)
		return "", fmt.Errorf("Failed to marshal value: %v", err)
	}
	return string(linksJSON), nil
}

func UnmarshalFromJSON(data string, target interface{}) error {
	err := json.Unmarshal([]byte(data), target)
	if err != nil {
		logger.Log.Errorf("Error parsing JSON: %v", err)
		return fmt.Errorf("Error parsing JSON: %v", err)
	}
	return nil
}

func ExecQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		logger.Log.Errorf("Error executing query: %v", err)
		return nil, fmt.Errorf("Error executing query: %v", err)
	}
	return rows, nil
}

func ExecUpdate(query string, args ...interface{}) (int64, error) {
	res, err := DB.Exec(query, args...)
	if err != nil {
		logger.Log.Errorf("Error executing update query: %v", err)
		return 0, fmt.Errorf("Error executing update query: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Errorf("Error checking rows affected: %v", err)
		return 0, fmt.Errorf("Error checking rows affected: %v", err)
	}
	return rowsAffected, nil
}

// BuildUpdateQuery is a utility function that handles the construction of an UPDATE query
// and prepares the corresponding arguments, including marshaling JSON data for special fields (like links and tags).
func BuildUpdateQuery(updatedData map[string]interface{}) (string, []interface{}, error) {
	var query string
	var args []interface{}

	// dynamically add fields to the query based on the available data in updatedData
	for key, value := range updatedData {
        if key == "created_on" {
            key = "creation_date"
        }
		// the following switch statement should work fine for all
		// items that are, or can be strings,
		// I feel like this may look stupid now, but will revisit
		// if needs changes. This allows for only awkward
		// datatypes, like the links, to be handled differently.
		switch key {
		case "links", "tags":
			jsonData, err := MarshalToJSON(value)
			if err != nil {
				return "", nil, fmt.Errorf("Error marshaling list data for key `%v`: %v", key, err)
			}
			query += fmt.Sprintf("%v = ?, ", key)
			args = append(args, string(jsonData))
		default:
			query += fmt.Sprintf("%v = ?, ", key)
			args = append(args, value)
		}
	}

	// continue formatting query
	// get rid of trailing space and comma
	query = query[:len(query)-2]
    // NOTICE we DO NOT add the `WHERE` clause here
	return query, args, nil
}
