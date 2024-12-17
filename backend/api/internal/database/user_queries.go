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

// a function to retrieve the username given a user id
//
// input:
//
//	id (int64) - the id
//
// output:
//
//	string - the username
//	error
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

// a function to retrieve the user id given a username
//
// input:
//
//	string - the username
//
// output:
//
//	id (int64) - the id
//	error
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

// a function to get all user data based on a username
//
// input:
//
//	username (string) - the username to query for
//
// output:
//
//	*types.User - the user retrieved
//	error
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

// a function to create a user in the database
//
// input:
//
//	user (*types.User) - the user to be created
//
// output:
//
//	error
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

// a function that deletes a user by their username
//
// input:
//
//	username (string) - the username to be deleted
//
// output:
//
//	int16 - status code
//	error
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

// a function that updates a user in the database given their username
//
// input:
//
//	username (string) - the username to be updated
//	updatedData (map[string]interface{}) - the updated data in JSON-like form
//
// output:
//
//	error
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

//	function to retrieve the usernames of a user's followers
//
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]string - the usernames
//	int - http status code
//	error
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

// function to retrieve the ids of a user's followers
//
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]int - the int ids
//	int - http status code
//	error
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
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]string - the usernames
//	int - http status code
//	error
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
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]int - the int ids
//	int - http status code
//	error
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

// helper function to use for retrieving
//
// input:
//
//	query (string) - the query used to get users data
//	userID (int) - the user's int
//
// output:
//
//	[]int - the ids of the retrieved users
//	int - httpcode
//	error
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

// helper function to use for retrieving
//
// input:
//
//	query (string) - the query used to get users data
//	userID (int) - the user's int
//
// output:
//
//	[]int - the ids of the retrieved users
//	int - httpcode
//	error
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

// function to create a follow from user to follow newFollow
//
// input:
//
//	user (string) - the username of the user requesting a follow
//	newFollow (string) - the username of the user to be followed
//
// output:
//
//	int - http code
//	error
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

// function to remove a follow from user to unfollow unfollow
//
// input:
//
//	user (string) - the username of the user requesting an unfollow
//	unfollow (string) - the username of the user to be unfollowed
//
// output:
//
//	int - http code
//	error
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
