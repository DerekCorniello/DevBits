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
	strId := context.Param("id")
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
		RespondWithError(context, http.StatusFound, fmt.Sprintf("Project with id '%v' not found", strId))
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

	username, err := database.GetUsernameById(newProj.Owner)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to verify project ownership: %v", err))
		return
	}

	if username == "" {
		RespondWithError(context, http.StatusNotFound, fmt.Sprintf("Failed to verify project ownership. User could not be found"))
		return
	}

	id, err := database.QueryCreateProject(&newProj)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to create project: %v", err))
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Project created successfully with id '%v'", id)})
}

func DeleteProject(context *gin.Context) {
	strId := context.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project ID: %v", err))
		return
	}

	code, err := database.QueryDeleteProject(id)
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
	context.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Project %v deleted.", id),
	})
}

func UpdateProjectInfo(context *gin.Context) {
	var updateData map[string]interface{}

	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to parse project id: %v", err))
		return
	}

	err = context.BindJSON(&updateData)
	if err != nil {
		RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Failed to update project: %v", err))
		return
	}

	existingProj, err := database.QueryProject(id)
	if existingProj == nil {
		RespondWithError(context, http.StatusNotFound, "Project not found")
		return
	}

	updatedData := make(map[string]interface{})

	// Iterate through the fields of the existing project and map the request data to those fields
	for key, value := range updateData {
		// use helper to check if the field exists in existingProj
		if IsFieldAllowed(existingProj, key) {
			updatedData[key] = value
		} else {
            RespondWithError(context, http.StatusBadRequest, fmt.Sprintf("Failed to update project: %v", err))
			return
		}
	}

	err = database.QueryUpdateProject(id, updatedData)
	if err != nil {
        RespondWithError(context, http.StatusInternalServerError, fmt.Sprintf("Error updating project: %v", err))
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": id})
}
