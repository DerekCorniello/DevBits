package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/api/internal/types"
)

// GetUsernameById retrieves the username for a given user ID.
func GetUsernameById(id int64) (string, error) {
	query := `SELECT username FROM Users WHERE id = ?;`

	row := DB.QueryRow(query, id)

	var retrievedUserName string
	err := row.Scan(&retrievedUserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return retrievedUserName, nil
}

// QueryUsername fetches the details of a user by their username.
func QueryUsername(username string) (*types.User, error) {
	query := `SELECT username, profile_pic, bio, links, creation_date FROM Users WHERE username = ?;`

	row := DB.QueryRow(query, username)

	var user types.User
	var linksJSON string
	err := row.Scan(&user.Username, &user.Picture, &user.Bio, &linksJSON, &user.CreationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse links JSON into a []string
	var links []string
	err = json.Unmarshal([]byte(linksJSON), &links)
	if err != nil {
		return nil, err
	}

	user.Links = links
	return &user, nil
}

// QueryCreateUser creates a new user in the database.
func QueryCreateUser(user *types.User) error {
	linksJSON, err := json.Marshal(user.Links)
	if err != nil {
		return fmt.Errorf("Failed to marshal links for user `%v`: %v", user.Username, err)
	}

	query := `INSERT INTO Users (username, profile_pic, bio, links)
	VALUES (?, ?, ?, ?);`

	res, err := DB.Exec(query, user.Username, user.Picture, user.Bio, string(linksJSON))
	if err != nil {
		return fmt.Errorf("Failed to create user `%v`: %v", user.Username, err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to ensure user was created: %v", err)
	}

	return nil
}

// QueryDeleteUser deletes a user from the database by username.
func QueryDeleteUser(username string) (int16, error) {
	query := `DELETE from Users WHERE username=?;`
	res, err := DB.Exec(query, username)
	if err != nil {
		return 400, fmt.Errorf("Failed to delete user `%v`: %v", username, err)
	}

	RowsAffected, err := res.RowsAffected()
	if RowsAffected == 0 {
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}

// QueryUpdateUser updates a user's information in the database.
func QueryUpdateUser(username string, updatedData map[string]interface{}) error {
    query := `UPDATE Users SET `
	var args []interface{}

	queryParams, args, err := BuildUpdateQuery(updatedData)
	if err != nil {
        return fmt.Errorf("Error building query: %v", err)
	}

	query += queryParams + " WHERE username = ?"
	args = append(args, username)

	rowsAffected, err := ExecUpdate(query, args...)
	if err != nil {
		return fmt.Errorf("Error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No user found with username `%s` to update", username)
	}

	return nil
}

// getUserIdByName retrieves the user ID by their username.
func getUserIdByName(username string) (int, error) {
	query := `SELECT id FROM Users WHERE username = ?`
	var userID int
	row := DB.QueryRow(query, username)
	err := row.Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("Error fetching user ID for username `%v`: %v", username, err)
	}
	return userID, nil
}

// QueryGetUsersFollowers retrieves the list of followers for a given username.
func QueryGetUsersFollowers(username string) ([]string, error) {
	userID, err := getUserIdByName(username)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT u.username 
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follower_id
        WHERE uf.follows_id = ?`

	return getUsersFollowingOrFollowers(query, userID)
}

// QueryGetUsersFollowing retrieves the list of users a given username is following.
func QueryGetUsersFollowing(username string) ([]string, error) {
	userID, err := getUserIdByName(username)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT u.username 
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follows_id
        WHERE uf.follower_id = ?`

	return getUsersFollowingOrFollowers(query, userID)
}

// getUsersFollowingOrFollowers is a shared function to fetch following or followers
func getUsersFollowingOrFollowers(query string, userID int) ([]string, error) {
	rows, err := ExecQuery(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
