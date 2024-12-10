package handlers

import (
	"fmt"
	"net/http"
    "strconv"

	"backend/api/internal/database"
	"backend/api/internal/logger"

	"github.com/gin-gonic/gin"
)

func GetProjectById(context *gin.Context) {
	strId := context.Param("id")
    id, err :=  strconv.Atoi(strId)
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
