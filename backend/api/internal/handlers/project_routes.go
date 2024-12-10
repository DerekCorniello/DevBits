package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/api/internal/database"
	"backend/api/internal/logger"
	"backend/api/internal/types"

	"github.com/gin-gonic/gin"
)

func GetProjectById(context *gin.Context) {
	strId := context.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Log.Infof("Failed to parse project ID: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": fmt.Sprintf("Failed to parse project ID: %v", err),
		})
		return
	}
	project, err := database.QueryProject(id)
	if err != nil {
		logger.Log.Infof("Failed to get project: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": fmt.Sprintf("Failed to fetch project: %v", err),
		})
		return
	}

	if project == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Project with id '%v' not found", strId),
		})
		return
	}

	context.JSON(http.StatusOK, project)
}

func CreateProject(context *gin.Context) {
	var newProj types.Project
	err := context.BindJSON(&newProj)

	if err != nil {
		logger.Log.Infof("Failed to bind to JSON: %v", err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

    username, err := database.GetUsernameById(newProj.Owner)
    if err != nil {
		logger.Log.Infof("Failed to verify project ownership: %v", err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
    }

    if username == "" {
		logger.Log.Info("Failed to verify project ownership. User could not be found")
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "User bound to project could not be found.",
		})
		return
    }

	id, err := database.QueryCreateProject(&newProj)
	if err != nil {
		logger.Log.Infof("Failed to create project: %v", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Project created successfully with id '%v'", id)})
}

func DeleteProject(context *gin.Context) {
	strId := context.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Log.Infof("Failed to parse project ID: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": fmt.Sprintf("Failed to parse project ID: %v", err),
		})
		return
	}

	code, err := database.QueryDeleteProject(id)
	if err != nil {
		logger.Log.Infof("Failed to delete project: %v", err)
		var httpCode int
		switch code {
		case 400:
			httpCode = http.StatusBadRequest
		case 404:
			httpCode = http.StatusNotFound
		default:
			httpCode = http.StatusInternalServerError
		}
		context.IndentedJSON(httpCode, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Project %v deleted.", id),
	})
}

func UpdateProjectInfo(context *gin.Context) {
	var updateData map[string]interface{}

	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		logger.Log.Infof("Failed to parse project id: %v", err)
	}

	err = context.BindJSON(&updateData)
	if err != nil {
		logger.Log.Infof("Failed to update user: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching user: %v", err)})
		return
	}

	existingProj, err := database.QueryProject(id)
	if existingProj == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
	}

	updatedData := make(map[string]interface{})

	// Iterate through the fields of the existing user and map the request data to those fields
	for key, value := range updateData {
		// use helper to check if the field exists in existingUser
		if IsFieldAllowed(existingProj, key) {
			updatedData[key] = value
		} else {
			logger.Log.Infof("Failed to update project: %v", err)
			context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' is not allowed to be updated", key)})
			return
		}
	}

	err = database.QueryUpdateProject(id, updatedData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating project: %v", err)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": id})
}
