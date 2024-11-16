package main

import (
    "backend/api/internal/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/v1/users", handlers.GetUsersV1)
	router.GET("/v1/users/:username", handlers.GetUserByUsernameV1)
	router.PUT("/v1/users/:username", handlers.UpdateUserInfoV1)
	router.DELETE("/v1/users/:username", handlers.DeleteUserV1)
	router.POST("/v1/users", handlers.CreateUserV1)

	router.Run("localhost:8080")
}
