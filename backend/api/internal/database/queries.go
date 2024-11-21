package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/api/internal/logger"
	"backend/api/internal/types"
)

func QueryUsername(username string) (*types.User, error) {
	query := `SELECT username, profile_pic, bio, links, creation_date FROM Users WHERE username = ?;`

	row := DB.QueryRow(query, username)

	var user types.User
	var linksJSON string

	err := row.Scan(&user.Username, &user.Picture, &user.Bio, &linksJSON, &user.CreationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Infof("No user found with username: %s", username)
			return nil, nil
		}
		logger.Log.Infof("Error fetching user: %v", err)
		return nil, err
	}

	// Parse links JSON into a []sql.NullString
	var links []string
	err = json.Unmarshal([]byte(linksJSON), &links)
	if err != nil {
		logger.Log.Infof("Error parsing links JSON: %v", err)
		return nil, err
	}

	return &user, nil
}

func QueryCreateUser(user *types.User) error {
	linksJSON, err := json.Marshal(user.Links)
	if err != nil {
		logger.Log.Errorf("Failed to marshal links for user `%v`: %v", user.Username, err)
		return fmt.Errorf("Failed to marshal links for user `%v`: %v", user.Username, err)
	}

	// backend should handle updating creation time to current time.
	query := `INSERT INTO Users (username, profile_pic, bio, links)
	VALUES (?, ?, ?, ?);`

	res, err := DB.Exec(query, user.Username, user.Picture, user.Bio, string(linksJSON))

	if err != nil {
		logger.Log.Errorf("Failed to create user `%v`: %v", user.Username, err)
		return fmt.Errorf("Failed to create user `%v`: %v", user.Username, err)
	}

	// we dont really need the last ID, but we can retrieve it to ensure
	// that we have something created
	lastId, err := res.LastInsertId()
	if err != nil {
		logger.Log.Errorf("Failed to ensure user was created: %v", err)
		return fmt.Errorf("Failed to ensure user was created: %v", err)
	}

	logger.Log.Infof("Created user %v with id `%v`", user.Username, lastId)
	return nil
}

func QueryDeleteUser(username string) (int16, error) {
	query := `DELETE from Users WHERE username=?;`
	res, err := DB.Exec(query, username)
	if err != nil {
		logger.Log.Errorf("Failed to delete user `%v`: %v", username, err)
		return 400, fmt.Errorf("Failed to delete user `%v`: %v", username, err)
	}

	RowsAffected, err := res.RowsAffected()
	if RowsAffected == 0 {
		logger.Log.Errorf("Deletion did not affect any records")
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		logger.Log.Errorf("Failed to fetch affected rows: %v", err)
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	logger.Log.Infof("Deleted user %v.", username)
	return 200, nil
}

func QueryUpdateUser(username string, updatedData map[string]interface{}) error {
	// lets use some string concat to ensure that we get good updates
	query := "UPDATE Users SET "
	var args []interface{}

	// dynamically add fields to the query based on the available data in updatedData
	for key, value := range updatedData {
		// the following switch statement should work fine for all
		// items that are, or can be strings,
		// I feel like this may look stupid now, but will revisit
		// if needs changes. This allows for only awkward
		// datatypes, like the links, to be handled differently.
		switch key {
		case "links":
			// parse links to JSON string
			linksJSON, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("Error marshaling links: %v", err)
			}
			query += "links = ?, "
			args = append(args, string(linksJSON))
		default:
			query += fmt.Sprintf("%v = ?, ", key)
			args = append(args, value)
		}
	}

	// continue formatting query
	// get rid of trailing space and comma
	query = query[:len(query)-2]
	query += " WHERE username = ?"
	args = append(args, username)

	res, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("Error executing update query: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No user found with username `%s` to update", username)
	}

	return nil
}
