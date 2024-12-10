package handlers

import (
	"fmt"
	"net/http"

	"backend/api/internal/database"
	"backend/api/internal/logger"
	"backend/api/internal/types"

	"github.com/gin-gonic/gin"
)


func GetUsernameById(context *gin.Context) {
	username := context.Param("username")

	user, err := database.QueryUsername(username)
	if err != nil {
		logger.Log.Infof("Failed to get user: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": fmt.Sprintf("Failed to fetch user: %v", err),
		})
		return
	}

	if user == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User with username '%v' not found", username),
		})
		return
	}

	context.JSON(http.StatusOK, user)
}

func GetUserByUsername(context *gin.Context) {
	username := context.Param("username")

	user, err := database.QueryUsername(username)
	if err != nil {
		logger.Log.Infof("Failed to get user: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": fmt.Sprintf("Failed to fetch user: %v", err),
		})
		return
	}

	if user == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User with username '%v' not found", username),
		})
		return
	}

	context.JSON(http.StatusOK, user)
}

func CreateUser(context *gin.Context) {
	var newUser types.User
	err := context.BindJSON(&newUser)

	if err != nil {
		logger.Log.Infof("Failed to bind to JSON: %v", err)
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	err = database.QueryCreateUser(&newUser)
	if err != nil {
		logger.Log.Infof("Failed to create user: %v", err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}
	context.IndentedJSON(http.StatusCreated, newUser)
}

func DeleteUser(context *gin.Context) {
	username := context.Param("username")
	code, err := database.QueryDeleteUser(username)
	if err != nil {
		logger.Log.Infof("Failed to delete user: %v", err)
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
		"message": fmt.Sprintf("User %v deleted.", username),
	})
}

func UpdateUserInfo(context *gin.Context) {
    // we dont want to create a whole new user, that is
    // why we dont use a user type here...
    // maybe could change later
	var updateData map[string]interface{}
	username := context.Param("username")

	// Bind the incoming JSON data to a map
    err := context.BindJSON(&updateData);
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	existingUser, err := database.QueryUsername(username)
	if err != nil {
		logger.Log.Infof("Failed to update user: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching user: %v", err)})
		return
	}

	if existingUser == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updatedData := make(map[string]interface{})

	// Iterate through the fields of the existing user and map the request data to those fields
	for key, value := range updateData {
		// use helper to check if the field exists in existingUser
		if IsFieldAllowed(existingUser, key) {
			updatedData[key] = value
		} else {
            logger.Log.Infof("Failed to update user: %v", err)
			context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' is not allowed to be updated", key)})
			return
		}
	}
    err = database.QueryUpdateUser(username, updatedData);
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating user: %v", err)})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": updatedData})
}

func GetUsersFollowers(context *gin.Context) {
	username := context.Param("username")

	followers, err := database.QueryGetUsersFollowers(username)
	if err != nil {
		logger.Log.Infof("Failed to obtain followers: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": fmt.Sprintf("Failed to fetch followers: %v", err),
		})
		return
	}

	if followers == nil || len(followers) == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("No followers for username '%v' found", username),
		})
		return
	}

	context.JSON(http.StatusOK, followers)
}

func GetUsersFollowing(context *gin.Context) {
	username := context.Param("username")

	followers, err := database.QueryGetUsersFollowing(username)
	if err != nil {
		logger.Log.Infof("Failed to obtain following: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": fmt.Sprintf("Failed to fetch following: %v", err),
		})
		return
	}

	if followers == nil || len(followers) == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("No followers for username '%v' found", username),
		})
		return
	}

	context.JSON(http.StatusOK, followers)
}
