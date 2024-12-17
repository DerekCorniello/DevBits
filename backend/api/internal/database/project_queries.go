package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"backend/api/internal/types"
)

// a function to get all project data based on a username
//
// input:
//
//	id (int) - the project id to query for
//
// output:
//
//	*types.Project - the user retrieved
//	error
func QueryProject(id int) (*types.Project, error) {
	query := `SELECT id, name, description, status, likes, links, tags, owner, creation_date FROM Projects WHERE id = ?;`
	row := DB.QueryRow(query, id)
	var project types.Project
	var linksJSON, tagsJSON string

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Status,
		&project.Likes,
		&linksJSON,
		&tagsJSON,
		&project.Owner,
		&project.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := UnmarshalFromJSON(linksJSON, &project.Links); err != nil {
		return nil, err
	}
	if err := UnmarshalFromJSON(tagsJSON, &project.Tags); err != nil {
		return nil, err
	}

	return &project, nil
}

// a function to create a project in the database
//
// input:
//
//	user (*types.Project) - the project to be created
//
// output:
//
//	 int64 - last project id created, for verifcation
//		error
func QueryCreateProject(proj *types.Project) (int64, error) {
	linksJSON, err := MarshalToJSON(proj.Links)
	if err != nil {
		return -1, err
	}

	tagsJSON, err := MarshalToJSON(proj.Tags)
	if err != nil {
		return -1, err
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Projects (name, description, status, links, tags, owner, creation_date)
              VALUES (?, ?, ?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, proj.Name, proj.Description, proj.Status, string(linksJSON), string(tagsJSON), proj.Owner, currentTime)
	if err != nil {
		return -1, fmt.Errorf("Failed to create project `%v`: %v", proj.Name, err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure proj was created: %v", err)
	}

	return lastId, nil
}

// a function that deletes a project by its id
//
// input:
//
//	id (int) - the project id of the project to be deleted
//
// output:
//
//	int16 - status code
//	error
func QueryDeleteProject(id int) (int16, error) {
	query := `DELETE from Projects WHERE id=?;`
	res, err := DB.Exec(query, id)
	if err != nil {
		return 400, fmt.Errorf("Failed to delete project `%v`: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}

// a function that updates a project in the database given the user id
//
// input:
//
//	id (int) - the project id of the project to be updated
//	updatedData (map[string]interface{}) - the updated data in JSON-like form
//
// output:
//
//	error
func QueryUpdateProject(id int, updatedData map[string]interface{}) error {
	query := `UPDATE Projects SET `
	var args []interface{}

	queryParams, args, err := BuildUpdateQuery(updatedData)
	if err != nil {
		return fmt.Errorf("Error building query: %v", err)
	}

	query += queryParams + " WHERE id = ?"
	args = append(args, id)

	rowsAffected, err := ExecUpdate(query, args...)
	if err != nil {
		return fmt.Errorf("Error executing update query: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No project found with id `%d` to update", id)
	}

	return nil
}

// function to retrieve the ids of a project's followers
//
// input:
//
//	projectID (int) - the project to retrieve
//
// output:
//
//	[]int - the int ids
//	int - http status code
//	error
func QueryGetProjectFollowers(projectID int) ([]int, int, error) {
	query := `
        SELECT u.id
        FROM Users u
        JOIN ProjectFollows pf ON u.id = pf.user_id
        WHERE pf.project_id = ?`

	return getProjectFollowersOrFollowing(query, projectID)
}

// function to retrieve the usernames of a project's followers
//
// input:
//
//	projectID (int) - the user to retrieve
//
// output:
//
//	[]string - the string usernames
//	int - http status code
//	error
func QueryGetProjectFollowersUsernames(projectID int) ([]string, int, error) {
	query := `
        SELECT u.username
        FROM Users u
        JOIN ProjectFollows pf ON u.id = pf.user_id
        WHERE pf.project_id = ?`

	return getProjectFollowersOrFollowingUsernames(query, projectID)
}

// function to retrieve the projects a user is following by id
//
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]int - the int project ids
//	int - http status code
//	error
func QueryGetProjectFollowing(username string) ([]int, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("Error fetching user id from username: %v", err)
	}

	query := `
        SELECT p.id
        FROM Projects p
        JOIN ProjectFollows pf ON p.id = pf.project_id
        WHERE pf.user_id = ?`

	return getProjectFollowersOrFollowing(query, userID)
}

// function to retrieve the projects a user is following by name
//
// input:
//
//	username (string) - the user to retrieve
//
// output:
//
//	[]string - the string project names
//	int - http status code
//	error
func QueryGetProjectFollowingNames(username string) ([]string, int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("Error fetching user id from username: %v", err)
	}

	query := `
        SELECT p.name
        FROM Projects p
        JOIN ProjectFollows pf ON p.id = pf.project_id
        WHERE pf.user_id = ?`

	return getProjectFollowersOrFollowingUsernames(query, userID)
}

// helper function to use for retrieving followers and following
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
func getProjectFollowersOrFollowing(query string, userID int) ([]int, int, error) {
	rows, err := ExecQuery(query, userID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var projectIDs []int
	for rows.Next() {
		var projectID int
		if err := rows.Scan(&projectID); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		projectIDs = append(projectIDs, projectID)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return projectIDs, http.StatusOK, nil
}

func getProjectFollowersOrFollowingUsernames(query string, userID int) ([]string, int, error) {
	rows, err := ExecQuery(query, userID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var projectNames []string
	for rows.Next() {
		var projectName string
		if err := rows.Scan(&projectName); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		projectNames = append(projectNames, projectName)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return projectNames, http.StatusOK, nil
}

func CreateNewProjectFollow(username string, projectID string) (int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred getting id for username: %v", username)
	}

	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred parsing project id: %v", projectID)
	}
	currFollowing, httpCode, err := QueryGetProjectFollowers(userID)
	if err != nil {
		return httpCode, fmt.Errorf("Cannot retrieve user's following list: %v", err)
	}

	if slices.Contains(currFollowing, intProjectID) {
		return http.StatusConflict, fmt.Errorf("User is already following this project")
	}

	query := `INSERT INTO ProjectFollows (user_id, project_id) VALUES (?, ?)`
	rowsAffected, err := ExecUpdate(query, userID, projectID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred adding project follow: %v", err)
	}

	if rowsAffected == 0 {
		return http.StatusInternalServerError, fmt.Errorf("Failed to add the follow relationship")
	}

	return http.StatusOK, nil
}

func RemoveProjectFollow(username string, projectID string) (int, error) {
	userID, err := GetUserIdByUsername(username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred getting id for username: %v", username)
	}

	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred parsing project id: %v", projectID)
	}

	currFollowing, httpCode, err := QueryGetProjectFollowers(userID)
	if err != nil {
		return httpCode, fmt.Errorf("Error retrieving user's following list: %v", err)
	}

	if !slices.Contains(currFollowing, intProjectID) {
		return http.StatusConflict, fmt.Errorf("User is not following this project")
	}

	query := `DELETE FROM ProjectFollows WHERE user_id = ? AND project_id = ?`
	rowsAffected, err := ExecUpdate(query, userID, projectID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred removing project follow: %v", err)
	}

	if rowsAffected == 0 {
		return http.StatusConflict, fmt.Errorf("No such follow relationship exists")
	}

	return http.StatusOK, nil
}
