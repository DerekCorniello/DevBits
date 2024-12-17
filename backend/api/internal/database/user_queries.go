package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"backend/api/internal/types"
)

// GetUsernameById retrieves the username associated with the given user ID.
//
// Parameters:
//   - id: The unique identifier of the user.
//
// Returns:
//   - string: The username if found.
//   - error: An error if the query fails. Returns nil for both if no user exists.
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

// GetUserIdByUsername retrieves the user ID associated with the given username.
//
// Parameters:
//   - username: The username to query.
//
// Returns:
//   - int: The user ID if found.
//   - error: An error if the query fails or the username does not exist.
func GetUserIdByUsername(username string) (int, error) {
	query := `SELECT id FROM Users WHERE username = ?`
	var userID int
	row := DB.QueryRow(query, username)
	err := row.Scan(&userID)
	if err != nil { // TODO: Is there a way this can return a 404 vs 500 error? this could be a 404 or 500, but we cannot tell from an err here
		return -1, fmt.Errorf("Error fetching user ID for username '%v' (this usually means username does not exist) : %v", username, err)
	}
	return userID, nil
}

// QueryUsername retrieves all user data for the given username.
//
// Parameters:
//   - username: The username to query.
//
// Returns:
//   - *types.User: The user details if found.
//   - error: An error if the query or data parsing fails.
func QueryUsername(username string) (*types.User, error) {
	query := `SELECT username, picture, bio, links, creation_date FROM Users WHERE username = ?;`

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

	var links []string
	err = json.Unmarshal([]byte(linksJSON), &links)
	if err != nil {
		return nil, err
	}

	user.Links = links
	return &user, nil
}

// QueryCreateUser creates a new user in the database.
//
// Parameters:
//   - user: The user data to insert.
//
// Returns:
//   - error: An error if the user creation fails.
func QueryCreateUser(user *types.User) error {
	linksJSON, err := json.Marshal(user.Links)
	if err != nil {
		return fmt.Errorf("Failed to marshal links for user '%v': %v", user.Username, err)
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Users (username, picture, bio, links, creation_date)
	VALUES (?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, user.Username, user.Picture, user.Bio, string(linksJSON), currentTime)
	if err != nil {
		return fmt.Errorf("Failed to create user '%v': %v", user.Username, err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Failed to ensure user was created: %v", err)
	}

	return nil
}

// QueryDeleteUser deletes a user by their username.
//
// Parameters:
//   - username: The username of the user to delete.
//
// Returns:
//   - int16: HTTP-like status code indicating the result.
//   - error: An error if the deletion fails.
func QueryDeleteUser(username string) (int16, error) {
	query := `DELETE from Users WHERE username=?;`
	res, err := DB.Exec(query, username)
	if err != nil {
		return 400, fmt.Errorf("Failed to delete user '%v': %v", username, err)
	}

	RowsAffected, err := res.RowsAffected()
	if RowsAffected == 0 {
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}

// QueryUpdateUser updates a user's details by their username.
//
// Parameters:
//   - username: The username of the user to update.
//   - updatedData: A map of fields to update and their new values.
//
// Returns:
//   - error: An error if the update fails or no user is found.
func QueryUpdateUser(username string, updatedData map[string]interface{}) error {

	newUsername, usernameExists := updatedData["username"]
	usernameStr, parseOk := newUsername.(string)

	// if there is a new username provided, ensure it is not empty
	if usernameExists && parseOk && usernameStr == "" {
		return fmt.Errorf("Updated username cannot be empty!")
	}

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
		return fmt.Errorf("No user found with username '%s' to update", username)
	}

	return nil
}

// QueryGetUsersFollowersUsernames retrieves the usernames of users who follow the specified user.
//
// Parameters:
//   - username: The username of the user.
//
// Returns:
//   - []string: A list of usernames of the followers.
//   - int: HTTP-like status code.
//   - error: An error if the query fails.
func QueryGetUsersFollowersUsernames(username string) ([]string, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query := `
        SELECT u.username
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follower_id
        WHERE uf.follows_id = ?`

	return getUsersFollowingOrFollowersUsernames(query, userID)
}

// function to retrieve the user ids of the users who follow the given user
//
// Parameters:
//  - username (string): the user to retrieve
//
// Returns:
//	- []int: a list of user ids of users who follow the specified user
//	- int: HTTP status code indicating the result of the operation
//	- error: any error encountered during the query
func QueryGetUsersFollowers(username string) ([]int, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query := `
        SELECT u.id 
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follower_id
        WHERE uf.follows_id = ?`

	return getUsersFollowingOrFollowers(query, userID)
}

// function to retrieve the usernames of the users who follow the given user
//
// Parameters:
//	- username (string): the user to retrieve
//
// Returns:
//	- []string: a list of usernames of users who follow the specified user
//	- int: HTTP status code indicating the result of the operation
//	- error: any error encountered during the query
func QueryGetUsersFollowingUsernames(username string) ([]string, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	query := `
        SELECT u.username 
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follows_id
        WHERE uf.follower_id = ?`

	return getUsersFollowingOrFollowersUsernames(query, userID)
}

// function to retrieve the ids of the users who follow the given user
//
// Parameters:
//	- username (string) - the user to retrieve
//
// Returns:
//	- []int: a list of user IDs of users who follow the specified user
//	- int: HTTP status code indicating the result of the operation
//	- error: any error encountered during the query
func QueryGetUsersFollowing(username string) ([]int, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	query := `
        SELECT u.id 
        FROM Users u
        JOIN UserFollows uf ON u.id = uf.follows_id
        WHERE uf.follower_id = ?`

	return getUsersFollowingOrFollowers(query, userID)
}

// helper function to retrieve the followers or followings of a user by their IDs
//
// Parameters:
//	- query (string): the SQL query to execute
//	- userID (int): the ID of the user to find follow data for
//
// Returns:
//	- []int: a list of user IDs for the followers or followings
//	- int: HTTP status code
//	- error: any error encountered during the query
func getUsersFollowingOrFollowers(query string, userID int) ([]int, int, error) {
	rows, err := ExecQuery(query, userID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var username int
		if err := rows.Scan(&username); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		users = append(users, username)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return users, http.StatusOK, nil
}

// helper function to retrieve the followers or followings of a user by their usernames
//
// Parameters:
//	- query (string): the SQL query to execute
//	- userID (int): the ID of the user to find follow data for
//
// Returns:
//	- []string: a list of usernames for the followers or followings
//	- int: HTTP status code
//	- error: any error encountered during the query
func getUsersFollowingOrFollowersUsernames(query string, userID int) ([]string, int, error) {
	rows, err := ExecQuery(query, userID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		users = append(users, username)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return users, http.StatusOK, nil
}

// function to create a follow relationship between two users
//
// Parameters:
//	- user (string): the username of the user initiating the follow
//	- newFollow (string): the username of the user to be followed
//
// Returns:
//	- int: HTTP status code
//	- error: any error encountered during the query
func CreateNewUserFollow(user string, newFollow string) (int, error) {
	userID, err := GetUserIdByUsername(user)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Cannot find user with username '%v'.", user)
	}

	newFollowID, err := GetUserIdByUsername(newFollow)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Cannot find user with username '%v'.", newFollow)
	}

	currFollowers, httpCode, err := QueryGetUsersFollowing(user)
	if err != nil {
		return httpCode, fmt.Errorf("Cannot retrieve user's following list: %v", err)
	}

	if slices.Contains(currFollowers, newFollowID) {
		return http.StatusConflict, fmt.Errorf("User '%v' is already being followed", newFollow)
	}

	query := `INSERT INTO UserFollows (follower_id, follows_id) VALUES (?, ?)`
	rowsAffected, err := ExecUpdate(query, userID, newFollowID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred adding follower: %v", err)
	}

	if rowsAffected == 0 {
		return http.StatusInternalServerError, fmt.Errorf("Failed to add the follow relationship")
	}

	return http.StatusOK, nil
}

// function to remove a follow relationship between two users
//
// Parameters:
//	- user (string): the username of the user initiating the unfollow
//	- unfollow (string): the username of the user to be unfollowed
//
// Returns:
//	- int: HTTP status code
//	- error: any error encountered during the query
func RemoveUserFollow(user string, unfollow string) (int, error) {
	userID, err := GetUserIdByUsername(user)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Cannot find user with username '%v'", user)
	}

	unfollowID, err := GetUserIdByUsername(unfollow)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("Cannot find user with username '%v'", unfollow)
	}

	currFollowers, httpCode, err := QueryGetUsersFollowing(user)
	if err != nil {
		return httpCode, fmt.Errorf("Error retrieving user's following list: %v", err)
	}

	if !slices.Contains(currFollowers, unfollowID) {
		return http.StatusConflict, fmt.Errorf("User '%v' is not being followed", unfollow)
	}

	query := `DELETE FROM UserFollows WHERE follower_id = ? AND follows_id = ?;`
	rowsAffected, err := ExecUpdate(query, userID, unfollowID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred removing follower: %v", err)
	}

	if rowsAffected == 0 {
		return http.StatusConflict, fmt.Errorf("No such follow relationship exists")
	}

	return http.StatusOK, nil
}
