package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/api/internal/database"
	"backend/api/internal/types"

	"github.com/gin-gonic/gin"
)

func GetProjectById(context *gin.Context) {
	strId := context.Param("project_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return
	}
	project, err := database.QueryProject(id)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch project: %v", err))
		return
	}

	if project == nil {
		RespondWithError(context, http.StatusNotFound, fmt.Sprintf("Project with id '%v' not found", strId))
		return
	}

	context.JSON(http.StatusOK, project)
}

func CreateProject(context *gin.Context) {
	var newProj types.Project
	err := context.BindJSON(&newProj)

	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to bind to JSON: %v", err))
		return
	}

	// verify the owner
	username, err := database.GetUsernameById(newProj.Owner)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify project ownership: %v", err))
		return
	}

	if username == "" {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify project ownership. User could not be found"))
		return
	}

	id, err := database.QueryCreateProject(&newProj)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to create project: %v", err))
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Project created successfully with id '%v'", id)})
}

func DeleteProject(context *gin.Context) {
	strId := context.Param("project_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project id: %v", err))
		return
	}

	code, err := database.QueryDeleteProject(id)
	// delete projects can return different errors...
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
		RespondWithError(context, httpCode, fmt.Sprintf("Failed to delete project: %v", err))
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Project %v deleted.", id),
	})
}
func UpdateProjectInfo(context *gin.Context) {
	var updateData map[string]interface{}

	// Parse project ID from the URL
	id, err := strconv.Atoi(context.Param("project_id"))
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project id: %v", err))
		return
	}

	// Parse incoming JSON
	err = context.BindJSON(&updateData)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse update data: %v", err))
		return
	}

	// Check if the project exists
	existingProj, err := database.QueryProject(id)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve project: %v", err))
		return
	}
	if existingProj == nil {
		RespondWithError(context, http.StatusNotFound, fmt.Sprintf("Project with id '%v' not found", id))
		return
	}

	// Validate new owner if provided in update data
	if newOwner, ok := updateData["owner"]; ok {
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

	// Filter and validate update fields
	updatedData := make(map[string]interface{})
	for key, value := range updateData {
		if IsFieldAllowed(existingProj, key) {
			updatedData[key] = value
		} else {
			RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Field '%v' is not allowed for updates", key))
			return
		}
	}

	// Update the project in the database
	err = database.QueryUpdateProject(id, updatedData)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Error updating project: %v", err))
		return
	}

	updatedProj, err := database.QueryProject(id)

	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Error validating updated project: %v", err))
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Project updated successfully",
		"project": updatedProj,
	})
}

func GetProjectFollowers(context *gin.Context) {
	projectId := context.Param("project_id")
	intProjectId, err := strconv.Atoi(projectId)

	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project id: %V", err))
	}

	followers, httpcode, err := database.QueryGetProjectFollowers(intProjectId)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to fetch followers: %v", err))
		return
	}

	context.JSON(http.StatusOK, followers)
}

func GetProjectFollowing(context *gin.Context) {
	username := context.Param("username")

	following, httpcode, err := database.QueryGetProjectFollowing(username)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to fetch following: %v", err))
		return
	}

	context.JSON(http.StatusOK, following)
}

func GetProjectFollowersUsernames(context *gin.Context) {
	projectId := context.Param("project_id")
	intProjectId, err := strconv.Atoi(projectId)

	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project id: %V", err))
	}

	followers, httpcode, err := database.QueryGetProjectFollowersUsernames(intProjectId)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to fetch followers: %v", err))
		return
	}

	context.JSON(http.StatusOK, followers)
}

func GetProjectFollowingUsernames(context *gin.Context) {
	username := context.Param("username")

	following, httpcode, err := database.QueryGetProjectFollowingUsernames(username)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to fetch following: %v", err))
		return
	}

	context.JSON(http.StatusOK, following)
}

func FollowProject(context *gin.Context) {
	username := context.Param("username")
	projectId := context.Param("project_id")

	httpcode, err := database.CreateNewProjectFollow(username, projectId)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to add follower: %v", err))
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%v now follows %v", username, projectId)})
}

func UnfollowProject(context *gin.Context) {
	username := context.Param("username")
	projectId := context.Param("project_id")

	httpcode, err := database.RemoveProjectFollow(username, projectId)
	if err != nil {
		RespondWithError(context, httpcode, fmt.Sprintf("Failed to remove follower: %v", err))
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%v unfollowed %v", username, projectId)})
}
