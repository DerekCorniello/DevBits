package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"backend/api/internal/types"
)

// QueryPosts retrieves a post by its ID from the database.
//
// Parameters:
//   - id: The unique identifier of the post to query.
//
// Returns:
//   - *types.Post: The post details if found.
//   - error: An error if the query fails. Returns nil for both if no post exists.
func QueryPost(id int) (*types.Post, error) {
	query := `SELECT id, user_id, project_id, content, likes, creation_date FROM Posts WHERE id = ?;`
	row := DB.QueryRow(query, id)
	var post types.Post

	err := row.Scan(
		&post.ID,
		&post.User,
		&post.Project,
		&post.Content,
		&post.Likes,
		&post.CreationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

// QueryCreatePost creates a new post in the database.
//
// Parameters:
//   - post: The post to be created, containing all necessary fields.
//
// Returns:
//   - int64: The ID of the newly created post.
//   - error: An error if the operation fails.
func QueryCreatePost(post *types.Post) (int64, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Posts (id, user_id, project_id, content, likes, creation_date) 
              VALUES (?, ?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, post.ID, post.User, post.Project, post.Content, post.Likes, currentTime)
	if err != nil {
		return -1, fmt.Errorf("Failed to create post: %v", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure post was created: %v", err)
	}

	return lastId, nil
}

// QueryDeletePost deletes a post by its ID.
//
// Parameters:
//   - id: The unique identifier of the post to delete.
//
// Returns:
//   - int16: http status code indicating the result of the operation.
//   - error: An error if the operation fails or no post is found.
func QueryDeletePost(id int) (int16, error) {
	query := `DELETE from Posts WHERE id=?;`
	res, err := DB.Exec(query, id)
	if err != nil {
		return 400, fmt.Errorf("Failed to delete post `%v`: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return 404, fmt.Errorf("Deletion did not affect any records")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}

// QueryUpdateProject updates an existing post in the database.
//
// Parameters:
//   - id: The unique identifier of the post to update.
//   - updatedData: A map containing the fields to update with their new values.
//
// Returns:
//   - error: An error if the operation fails or no post is found.
func QueryUpdatePost(id int, updatedData map[string]interface{}) error {
	query := `UPDATE Posts SET `
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
		return fmt.Errorf("No post found with id `%d` to update", id)
	}

	return nil
}

// QueryPostsByUserId retrieves a set of posts by its owning user id from the database.
//
// Parameters:
//   - id: The unique identifier of the user to query.
//
// Returns:
//   - []types.Post: The post details if found.
//   - error: An error if the query fails. Returns nil for both if no post exists.
func QueryPostsByUserId(userId int) ([]types.Post, int, error) {
	query := `SELECT id, user_id, project_id, content, likes, creation_date FROM Posts WHERE user_id = ?;`

	rows, err := DB.Query(query, userId)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var posts []types.Post

	for rows.Next() {
		var post types.Post
		err := rows.Scan(
			&post.ID,
			&post.User,
			&post.Project,
			&post.Content,
			&post.Likes,
			&post.CreationDate,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return []types.Post{}, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
        posts = append(posts, post)
	}

	return posts, http.StatusOK, nil
}

// QueryPostsByProjectId retrieves a set of posts by its owning project id from the database.
//
// Parameters:
//   - id: The unique identifier of the project to query.
//
// Returns:
//   - *types.Post: The post details if found.
//   - error: An error if the query fails. Returns nil for both if no post exists.
func QueryPostsByProjectId(projId int) ([]types.Post, int, error) {
	query := `SELECT id, user_id, project_id, content, likes, creation_date FROM Posts WHERE project_id = ?;`

	rows, err := DB.Query(query, projId)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()

	var posts []types.Post = []types.Post{}

	for rows.Next() {
		var post types.Post
		err := rows.Scan(
			&post.ID,
			&post.User,
			&post.Project,
			&post.Content,
			&post.Likes,
			&post.CreationDate,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return []types.Post{}, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
        posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}
