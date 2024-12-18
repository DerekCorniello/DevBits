package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/api/internal/database"
	"backend/api/internal/types"

	"github.com/gin-gonic/gin"
)

// GetPostById handles GET requests to retrieve project information by its ID.
// It expects the `post_id` parameter in the URL and does not require a request body.
// Returns:
// - 400 Bad Request if the ID is invalid.
// - 404 Not Found if the post does not exist.
// - 500 Internal Server Error if the database query fails.
// On success, responds with a 200 OK status and the post details in JSON format.
func GetPostById(context *gin.Context) {
	strId := context.Param("post_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}
	post, err := database.QueryPost(id)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch post: %v", err))
		return
	}

	if post == nil {
		RespondWithError(context, http.StatusNotFound, fmt.Sprintf("Post with id '%v' not found", strId))
		return
	}

	context.JSON(http.StatusOK, post)
}

// CreatePost handles POST requests to create a new post
// It expects a JSON payload that can be bound to a `types.Post` object.
// Validates the provided owner's ID and ensures the user and project exist.
// Returns:
// - 400 Bad Request if the JSON payload is invalid or the owner/project cannot be verified.
// - 500 Internal Server Error if there is a database error.
// On success, responds with a 201 Created status and the new post ID in JSON format.
func CreatePost(context *gin.Context) {
	var newPost types.Post
	err := context.BindJSON(&newPost)

	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to bind to JSON: %v", err))
		return
	}

	// verify the owner
	username, err := database.GetUsernameById(newPost.User)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify post ownership: %v", err))
		return
	}

	if username == "" {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify post ownership. User could not be found"))
		return
	}

	// verify the project
	project, err := database.QueryPost(int(newPost.User))
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify post ownership: %v", err))
		return
	}

	if project == nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify post ownership. Ownin project could not be found"))
		return
	}

	id, err := database.QueryCreatePost(&newPost)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to create project: %v", err))
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Post created successfully with id '%v'", id)})
}

// DeletePost handles DELETE requests to delete a post.
// It expects the `post_id` parameter in the URL.
// Returns:
// - 400 Bad Request if the post_id is invalid.
// - 404 Not Found if no post is found with the given id.
// - 500 Internal Server Error if a database query fails.
// On success, responds with a 200 OK status and a message confirming the post deletion.
func DeletePost(context *gin.Context) {
	strId := context.Param("post_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse post id: %v", err))
		return
	}

	code, err := database.QueryDeletePost(id)
	// delete posts can return different errors...
	if err != nil {
		var httpCode int
		switch code {
		case 400:
			httpCode = http.StatusBadRequest
		case 404:
			httpCode = http.StatusNotFound
		default:
			httpCode = http.StatusInternalServerError
		}
		RespondWithError(context, httpCode, fmt.Sprintf("Failed to delete post: %v", err))
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Post %v deleted.", id),
	})
}

// UpdatePostInfo handles PATCH requests to update post information.
// It expects the `post_id` parameter in the URL and a JSON payload with update fields.
// Validates the post ID, checks for the existence of the post, and ensures the fields being updated are allowed.
// Returns:
// - 400 Bad Request for invalid input or disallowed fields.
// - 404 Not Found if the post does not exist.
// - 500 Internal Server Error for database errors.
// On success, responds with a 200 OK status and the updated post details in JSON format.
func UpdatePostInfo(context *gin.Context) {
	var updateData map[string]interface{}

	// Parse post ID from the URL
	id, err := strconv.Atoi(context.Param("post_id"))
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse post id: %v", err))
		return
	}

	// Parse incoming JSON
	err = context.BindJSON(&updateData)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse update data: %v", err))
		return
	}

	// Check if the post exists
	existingPost, err := database.QueryPost(id)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve post: %v", err))
		return
	}
	if existingPost == nil {
		RespondWithError(context, http.StatusNotFound, fmt.Sprintf("Post with id '%v' not found", id))
		return
	}

	// Validate new owner if provided in update data
	if newOwner, ok := updateData["user"]; ok {
		ownerID, ok := newOwner.(float64) // Assuming JSON numbers are decoded as float64
		if !ok {
			RespondWithError(context, http.StatusBadRequest, "Invalid owner id format")
			return
		}
		username, err := database.GetUsernameById(int64(ownerID))
		if err != nil || username == "" {
			RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Invalid owner id: %v", ownerID))
			return
		}
	}

	// Validate new owner if provided in update data
	if newProject, ok := updateData["project"]; ok {
		projectID, ok := newProject.(float64) // Assuming JSON numbers are decoded as float64
		if !ok {
			RespondWithError(context, http.StatusBadRequest, "Invalid project id format")
			return
		}
		existingProject, err := database.QueryProject(int(projectID))
		if err != nil || existingProject == nil {
			RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Invalid project id: %v", projectID))
			return
		}
	}

	// Filter and validate update fields
	updatedData := make(map[string]interface{})
	for key, value := range updateData {
		if IsFieldAllowed(existingPost, key) {
			updatedData[key] = value
		} else {
			RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Field '%v' is not allowed for updates", key))
			return
		}
	}

	// Update the post in the database
	err = database.QueryUpdatePost(id, updatedData)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Error updating post: %v", err))
		return
	}

	updatedPost, err := database.QueryPost(id)

	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Error validating updated post: %v", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    updatedPost,
	})
}
