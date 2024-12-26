package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
    "strconv"

	"backend/api/internal/types"
)

// QueryComment retrieves a comment by its ID from the database.
//
// Parameters:
//   - id: The unique identifier of the comment to query.
//
// Returns:
//   - *types.Comment: The comment details if found.
//   - error: An error if the query fails. Returns nil for both if no comment exists.
func QueryComment(id int) (*types.Comment, error) {
	query := `SELECT id, user_id, content, likes, creation_date, parent_comment_id FROM Comments WHERE id = ?;`
	row := DB.QueryRow(query, id)
	var comment types.Comment

	err := row.Scan(
		&comment.ID,
		&comment.User,
		&comment.Content,
		&comment.Likes,
		&comment.CreationDate,
		&comment.ParentComment,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &comment, nil
}

// QueryCommentsByUserId retrieves a set of comments by its owning user id from the database.
//
// Parameters:
//   - id: The unique identifier of the user to query.
//
// Returns:
//   - []types.Post: The post details if found.
//   - error: An error if the query fails. Returns nil for both if no comments exists.
func QueryCommentsByUserId(userId int) ([]types.Comment, int, error) {
	query := `
            SELECT 
                c.id AS comment_id,
                c.user_id,
                c.content,
                c.likes,
                c.creation_date,
                c.parent_comment_id
            FROM Comments c
            JOIN PostComments pc ON c.id = pc.comment_id
            WHERE c.user_id = ?;
    `

	postRows, err := DB.Query(query, userId)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer postRows.Close()

	query = `
            SELECT 
                c.id AS comment_id,
                c.user_id,
                c.content,
                (SELECT COUNT(*) FROM CommentLikes cl WHERE cl.comment_id = c.id) AS likes,
                c.creation_date,
                c.parent_comment_id
            FROM Comments c
            JOIN ProjectComments pc ON c.id = pc.comment_id
            WHERE pc.project_id = ?;
    `
	projRows, err := DB.Query(query, userId)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer projRows.Close()
	var comments []types.Comment

	for projRows.Next() {
		var comment types.Comment
		err := projRows.Scan(
			&comment.ID,
			&comment.User,
			&comment.Content,
			&comment.Likes,
			&comment.CreationDate,
			&comment.ParentComment,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, comment)
	}
	for postRows.Next() {
		var comment types.Comment
		err := postRows.Scan(
			&comment.ID,
			&comment.User,
			&comment.Content,
			&comment.Likes,
			&comment.CreationDate,
			&comment.ParentComment,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, comment)
	}

	return comments, http.StatusOK, nil
}

// QueryCommentsByProjectId retrieves a comment by its project ID from the database.
//
// Parameters:
//   - id: The unique identifier of the project to query.
//
// Returns:
//   - *types.Comment: The comment details if found.
//   - error: An error if the query fails. Returns nil for both if no comment exists.
func QueryCommentsByProjectId(id int) ([]types.Comment, int, error) {
	query := `
            SELECT 
                c.id AS comment_id,
                c.user_id,
                c.content,
                c.likes,
                c.creation_date,
                c.parent_comment_id
            FROM Comments c
            JOIN ProjectComments pc ON c.id = pc.comment_id
            WHERE pc.project_id = ?;
    `
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()
	var comments []types.Comment

	for rows.Next() {
		var comment types.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.User,
			&comment.Content,
			&comment.Likes,
			&comment.CreationDate,
			&comment.ParentComment,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, comment)
	}
	return comments, http.StatusOK, nil
}

// QueryCommentsByPostId retrieves a comment by its post ID from the database.
//
// Parameters:
//   - id: The unique identifier of the post to query.
//
// Returns:
//   - *types.Comment: The comment details if found.
//   - error: An error if the query fails. Returns nil for both if no comment exists.
func QueryCommentsByPostId(id int) ([]types.Comment, int, error) {
	query := `
            SELECT 
                c.id AS comment_id,
                c.user_id,
                c.content,
                c.likes,
                c.creation_date,
                c.parent_comment_id
            FROM Comments c
            JOIN PostComments pc ON c.id = pc.comment_id
            WHERE pc.post_id = ?;
    `
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()
	var comments []types.Comment

	for rows.Next() {
		var comment types.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.User,
			&comment.Content,
			&comment.Likes,
			&comment.CreationDate,
			&comment.ParentComment,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, comment)
	}
	return comments, http.StatusOK, nil
}

// QueryCommentsByCommentId retrieves a comment by its comment ID from the database.
//
// Parameters:
//   - id: The unique identifier of the comment to query.
//
// Returns:
//   - *types.Comment: The comment details if found.
//   - error: An error if the query fails. Returns nil for both if no comment exists.
func QueryCommentsByCommentId(id int) ([]types.Comment, int, error) {
	query := `
            SELECT 
                c.id AS comment_id,
                c.user_id,
                c.content,
                c.likes,
                c.creation_date,
                c.parent_comment_id
            FROM Comments c
            WHERE c.parent_comment_id = ?;
    `
	rows, err := DB.Query(query, id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}
	defer rows.Close()
	var comments []types.Comment

	for rows.Next() {
		var comment types.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.User,
			&comment.Content,
			&comment.Likes,
			&comment.CreationDate,
			&comment.ParentComment,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, http.StatusOK, nil
			}
			return nil, http.StatusInternalServerError, err
		}
		comments = append(comments, comment)
	}
	return comments, http.StatusOK, nil
}

// QueryCreateCommentOnPost creates a new comment on a post in the database.
//
// Parameters:
//   - commment: The comment to be created, containing all necessary fields.
//   - postId: The id of the post for the comment to be added to
//
// Returns:
//   - int64: The ID of the newly created comment.
//   - error: An error if the operation fails.
func QueryCreateCommentOnPost(comment types.Comment, postId int) (int64, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Comments (user_id, content, parent_comment_id, likes, creation_date) 
              VALUES (?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, comment.User, comment.Content, comment.ParentComment, 0, currentTime)

	if err != nil {
		return -1, fmt.Errorf("Failed to create comment: %v", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure post was created: %v", err)
	}

	query = `INSERT INTO PostComments (user_id, post_id, comment_id)
             VALUES (?, ?, ?)`

	res, err = DB.Exec(query, comment.User, postId, lastId)
	if err != nil {
		return -1, fmt.Errorf("Failed to link comment to post: %v", err)
	}

	return lastId, nil
}

// QueryCreateCommentOnProject creates a new comment on a project in the database.
//
// Parameters:
//   - commment: The comment to be created, containing all necessary fields.
//   - projectId: The id of the project for the comment to be added to
//
// Returns:
//   - int64: The ID of the newly created comment.
//   - error: An error if the operation fails.
func QueryCreateCommentOnProject(comment types.Comment, projectId int) (int64, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Comments (user_id, content, parent_comment_id, likes, creation_date) 
              VALUES (?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, comment.User, comment.Content, comment.ParentComment, 0, currentTime)

	if err != nil {
		return -1, fmt.Errorf("Failed to create comment: %v", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure project was created: %v", err)
	}

	query = `INSERT INTO ProjectComments (user_id, project_id, comment_id)
             VALUES (?, ?, ?)`

	res, err = DB.Exec(query, comment.User, projectId, lastId)
	if err != nil {
		return -1, fmt.Errorf("Failed to link comment to project: %v", err)
	}

	return lastId, nil
}

// QueryCreateCommentOnComment creates a new comment on a comment in the database.
//
// Parameters:
//   - commment: The comment to be created, containing all necessary fields.
//   - commentId: The id of the comment for the comment to be added to
//
// Returns:
//   - int64: The ID of the newly created comment.
//   - error: An error if the operation fails.
func QueryCreateCommentOnComment(comment types.Comment, commentId int) (int64, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO Comments (user_id, content, parent_comment_id, likes, creation_date) 
              VALUES (?, ?, ?, ?, ?);`

	res, err := DB.Exec(query, comment.User, comment.Content, commentId, 0, currentTime)

	if err != nil {
		return -1, fmt.Errorf("Failed to create comment: %v", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to ensure comment was created: %v", err)
	}

	return lastId, nil
}

// QueryDeleteComment soft deletes a comment
//
// Parameters:
//   - id: The id of the comment to be deleted
//
// Returns:
//   - int16: http status code
//   - error: An error if the operation fails.
func QueryDeleteComment(id int) (int16, error) {
	_, err := DB.Exec(`UPDATE PostComments SET user_id = -1 WHERE comment_id = ?`, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to update PostComments for deleted comment: %v", err)
	}

	_, err = DB.Exec(`UPDATE ProjectComments SET user_id = -1 WHERE comment_id = ?`, id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to update ProjectComments for deleted comment: %v", err)
	}

	query := `UPDATE Comments SET user_id = -1, content = "This comment was deleted.", likes = 0 WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("Failed to soft delete comment `%v`: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return http.StatusNotFound, fmt.Errorf("Comment not found or already marked as deleted")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return http.StatusOK, nil
}

// QueryUpdateCommentContent updates comment's content
//
// Parameters:
//   - id: The id of the comment to be deleted
//   - newContent: the updated content
//
// Returns:
//   - int16: http status code
//   - error: An error if the operation fails.
func QueryUpdateCommentContent(id int, newContent string) (int16, error) {
	// get comment creation time to validate time diff
	var creationDate string
	query := `SELECT creation_date FROM Comments WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&creationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 404, fmt.Errorf("Comment not found")
		}
		return 500, fmt.Errorf("Failed to fetch comment creation date: %v", err)
	}

	createdAt, err := time.Parse("2006-01-02T15:04:05Z", creationDate)
	if err != nil {
		return 500, fmt.Errorf("Failed to parse creation date: %v", err)
	}

	now := time.Now()
	// now check if the comment is older than 2 minutes
	if now.Sub(createdAt) > 2*time.Minute {
		return 400, fmt.Errorf("Cannot update comment. More than 2 minutes have passed since posting.")
	}

	query = `UPDATE Comments SET content = ? WHERE id = ?`
	res, err := DB.Exec(query, newContent, id)
	if err != nil {
		return 500, fmt.Errorf("Failed to update comment content: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return 404, fmt.Errorf("Comment not found or no changes made")
	} else if err != nil {
		return 500, fmt.Errorf("Failed to fetch affected rows: %v", err)
	}

	return 200, nil
}


// CreateCommentLike creates a like relationship between a user and a comment.
//
// Parameters:
//   - username: The username of the user creating the like.
//   - strCommentID: The ID of the comment to like (as a string, converted internally).
//
// Returns:
//   - int: HTTP-like status code indicating the result of the operation.
//   - error: An error if the operation fails or the user is not liking the comment.
func CreateCommentLike(username string, strCommentId string) (int, error) {
	// get user ID from username, implicitly checks if user exists
	user_id, err := GetUserIdByUsername(username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred getting id for username: %v", err)
	}

	// parse comment ID
	commentId, err := strconv.Atoi(strCommentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred parsing user id: %v", err)
	}

	// verify comment exists
	_, err = QueryComment(commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred verifying the comment exists: %v", err)
	}

	// check if the like already exists
	var exists bool
	query := `SELECT EXISTS (
                 SELECT 1 FROM CommentLikes WHERE user_id = ? AND post_id = ?
              )`
	err = DB.QueryRow(query, user_id, commentId).Scan(&exists)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred checking like existence: %v", err)
	}
	if exists {
		// like already exists, but we return success to keep it idempotent
		return http.StatusCreated, nil
	}

	// insert the like
	insertQuery := `INSERT INTO CommentLikes (user_id, post_id) VALUES (?, ?)`
	_, err = DB.Exec(insertQuery, user_id, commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to insert comment like: %v", err)
	}

	// update the likes column
	updateQuery := `UPDATE Comments SET likes = likes + 1 WHERE id = ?`
	_, err = DB.Exec(updateQuery, commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to update likes count: %v", err)
	}

	return http.StatusCreated, nil
}

// RemoveCommentLike deletes a like relationship between a user and a comment.
//
// Parameters:
//   - username: The username of the user removing the like.
//   - strPostId: The ID of the comment to unlike (as a string, converted internally).
//
// Returns:
//   - int: HTTP-like status code indicating the result of the operation.
//   - error: An error if the operation fails or the user is not liking the comment.
func RemoveCommentLike(username string, strCommentId string) (int, error) {
	// get user ID
	user_id, err := GetUserIdByUsername(username)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred getting id for username: %v", err)
	}

	// parse post ID
	commentId, err := strconv.Atoi(strCommentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred parsing username id: %v", err)
	}

	// verify post exists
	_, err = QueryPost(commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("An error occurred verifying the comment exists: %v", err)
	}

	// perform the delete operation
	deleteQuery := `DELETE FROM CommentLikes WHERE user_id = ? AND post_id = ?`
	result, err := DB.Exec(deleteQuery, user_id, commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to delete comment like: %v", err)
	}

	// check if any rows were actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		// if no rows were deleted, return success to keep idempotency
		return http.StatusNoContent, nil
	}

	// update the likes column
	updateQuery := `UPDATE Comments SET likes = likes - 1 WHERE id = ?`
	_, err = DB.Exec(updateQuery, commentId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Failed to update likes count: %v", err)
	}

	return http.StatusOK, nil
}
