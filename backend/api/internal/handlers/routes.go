package handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"backend/api/internal/database"
	"backend/api/internal/logger"
	"backend/api/internal/types"

	"github.com/gin-gonic/gin"
)

// Here are all the routes

func GetUserByUsername(context *gin.Context) {
	username := context.Param("username")

	user, err := database.QueryUsername(username)
	if err != nil {
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

func isFieldAllowed(existingUser interface{}, fieldName string) bool {
	// existingUser should be a pointer to the struct, so get the type of the struct
	val := reflect.ValueOf(existingUser)

	// data insurance
	// If it's a pointer, dereference it to get the value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}

	// Loop through all fields of the struct
	for i := 0; i < val.NumField(); i++ {
		// Get the field and its JSON tag
		field := val.Type().Field(i)
		jsonTag := field.Tag.Get("json")

		// If the JSON tag matches the fieldName, return true
		if strings.ToLower(jsonTag) == strings.ToLower(fieldName) {
			return true
		}
	}

	return false
}

func UpdateUserInfo(context *gin.Context) {
	var updateData map[string]interface{}
	username := context.Param("username")

	// Bind the incoming JSON data to a map
	if err := context.BindJSON(&updateData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	existingUser, err := database.QueryUsername(username)
	if err != nil {
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
		if isFieldAllowed(existingUser, key) {
			updatedData[key] = value
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' is not allowed to be updated", key)})
			return
		}
	}

	if err := database.QueryUpdateUser(username, updatedData); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating user: %v", err)})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": updatedData})
}
