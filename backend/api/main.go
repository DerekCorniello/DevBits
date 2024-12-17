package main

import (
	"log"
	"os"

	"backend/api/internal/database"
	"backend/api/internal/handlers"
	"backend/api/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const DEBUG bool = true

func HealthCheck(context *gin.Context) {
	context.JSON(200, gin.H{"message": "API is running!"})
}

func main() {
	logger.InitLogger()

	router := gin.Default()
	router.HandleMethodNotAllowed = true

	// Apply CORS middleware to the router
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"}, // Add your frontend URL (React Native or Web app)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies or authentication headers
	}))

	router.GET("/health", HealthCheck)

	router.GET("/users/:username", handlers.GetUserByUsername)
	router.POST("/users", handlers.CreateUser)
	router.PUT("/users/:username", handlers.UpdateUserInfo)
	router.DELETE("/users/:username", handlers.DeleteUser)

	router.GET("/users/:username/followers", handlers.GetUsersFollowers)
	router.GET("/users/:username/follows", handlers.GetUsersFollowing)
	router.GET("/users/:username/followers/usernames", handlers.GetUsersFollowersUsernames)
	router.GET("/users/:username/follows/usernames", handlers.GetUsersFollowingUsernames)

    router.POST("/users/:username/follow/:new_follow", handlers.FollowUser)
    router.POST("/users/:username/unfollow/:unfollow", handlers.UnfollowUser)

    router.GET("/projects/:project_id", handlers.GetProjectById)
    router.POST("/projects", handlers.CreateProject)
	router.PUT("/projects/:project_id", handlers.UpdateProjectInfo)
    router.DELETE("/projects/:project_id", handlers.DeleteProject)

	router.GET("/projects/:project_id/followers", handlers.GetProjectFollowers)
    router.GET("/projects/follows/:username", handlers.GetProjectFollowing)
	router.GET("/projects/:project_id/followers/usernames", handlers.GetProjectFollowersUsernames)
    router.GET("/projects/follows/:username/names", handlers.GetProjectFollowingNames)

    router.POST("/projects/:username/follow/:project_id", handlers.FollowProject)
    router.POST("/projects/:username/unfollow/:project_id", handlers.UnfollowProject)


	var dbinfo, dbtype string
	if DEBUG {
		dbinfo = "./api/internal/database/dev.sqlite3"
		dbtype = "sqlite3"
	} else {
		dbinfo = os.Getenv("DB_INFO")
		dbtype = os.Getenv("DB_TYPE")
		if dbinfo == "" {
			log.Fatalln("FATAL: debug mode is false and 'DB_INFO' doesn't exist!")
		}
		if dbtype == "" {
			log.Fatalln("FATAL: debug mode is false and 'DB_TYPE' doesn't exist!")
		}
	}
	database.Connect(dbinfo, dbtype)

	router.Run("localhost:8080")
}
