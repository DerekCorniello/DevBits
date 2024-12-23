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

	// Users
	router.GET("/users/:username", handlers.GetUserByUsername)
	router.POST("/users", handlers.CreateUser)
	router.PUT("/users/:username", handlers.UpdateUserInfo)
	router.DELETE("/users/:username", handlers.DeleteUser)

	router.GET("/users/:username/followers", handlers.GetUsersFollowers)
	router.GET("/users/:username/following", handlers.GetUsersFollowing)
	router.GET("/users/:username/followers/usernames", handlers.GetUsersFollowersUsernames)
	router.GET("/users/:username/following/usernames", handlers.GetUsersFollowingUsernames)

	router.POST("/users/:username/follow/:target_user", handlers.FollowUser)
	router.POST("/users/:username/unfollow/:target_user", handlers.UnfollowUser)

	// Projects
	router.GET("/projects/:project_id", handlers.GetProjectById)
	router.POST("/projects", handlers.CreateProject)
	router.PUT("/projects/:project_id", handlers.UpdateProjectInfo)
	router.DELETE("/projects/:project_id", handlers.DeleteProject)

	router.GET("/users/:user_id/projects", handlers.GetProjectsByUserId)

	router.GET("/projects/:project_id/followers", handlers.GetProjectFollowers)
	router.GET("/users/:username/projects/following", handlers.GetProjectFollowing)
	router.GET("/projects/:project_id/followers/usernames", handlers.GetProjectFollowersUsernames)
	router.GET("/users/:username/projects/following/names", handlers.GetProjectFollowingNames)

	router.POST("/users/:username/follow/project/:project_id", handlers.FollowProject)
	router.POST("/users/:username/unfollow/project/:project_id", handlers.UnfollowProject)

	// Posts
	router.GET("/posts/:post_id", handlers.GetPostById)
	router.POST("/posts", handlers.CreatePost)
	router.PUT("/posts/:post_id", handlers.UpdatePostInfo)
	router.DELETE("/posts/:post_id", handlers.DeletePost)

	router.GET("/users/:user_id/posts", handlers.GetPostsByUserId)
	router.GET("/projects/:project_id/posts", handlers.GetPostsByProjectId)

	// Comments
	router.POST("/posts/:post_id/comments", handlers.CreateCommentOnPost)
	router.POST("/projects/:project_id/comments", handlers.CreateCommentOnProject)
	router.POST("/comments/:comment_id/reply", handlers.CreateCommentOnComment)
	router.GET("/comments/:comment_id", handlers.GetCommentById)
	router.PUT("/comments/:comment_id", handlers.UpdateCommentContent)
	router.DELETE("/comments/:comment_id", handlers.DeleteComment)

	router.GET("/users/:user_id/comments", handlers.GetCommentsByUserId)
	router.GET("/posts/:post_id/comments", handlers.GetCommentsByPostId)
	router.GET("/projects/:project_id/comments", handlers.GetCommentsByProjectId)
	router.GET("/comments/:comment_id/replies", handlers.GetCommentsByCommentId)

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
