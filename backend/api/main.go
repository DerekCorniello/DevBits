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
	if DEBUG {
		gin.SetMode(gin.DebugMode)
	}
	log.SetOutput(os.Stdout)
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
	router.GET("/projects/by-user/:user_id", handlers.GetProjectsByUserId)

	router.GET("/projects/:project_id/followers", handlers.GetProjectFollowers)
	router.GET("/projects/follows/:username", handlers.GetProjectFollowing)
	router.GET("/projects/:project_id/followers/usernames", handlers.GetProjectFollowersUsernames)
	router.GET("/projects/follows/:username/names", handlers.GetProjectFollowingNames)

	router.POST("/projects/:username/follow/:project_id", handlers.FollowProject)
	router.POST("/projects/:username/unfollow/:project_id", handlers.UnfollowProject)

	router.POST("/projects/:username/likes/:project_id", handlers.LikeProject)
	router.POST("/projects/:username/unlikes/:project_id", handlers.UnlikeProject)
	router.GET("/projects/does-like/:username/:project_id", handlers.IsProjectLiked)

	router.GET("/posts/:post_id", handlers.GetPostById)
	router.POST("/posts", handlers.CreatePost)
	router.PUT("/posts/:post_id", handlers.UpdatePostInfo)
	router.DELETE("/posts/:post_id", handlers.DeletePost)

	router.GET("/posts/by-user/:user_id", handlers.GetPostsByUserId)
	router.GET("/posts/by-project/:project_id", handlers.GetPostsByProjectId)

	router.POST("/posts/:username/likes/:post_id", handlers.LikePost)
	router.POST("/posts/:username/unlikes/:post_id", handlers.UnlikePost)
	router.GET("/posts/does-like/:username/:post_id", handlers.IsPostLiked)

	router.POST("/comments/for-post/:post_id", handlers.CreateCommentOnPost)
	router.POST("/comments/for-project/:project_id", handlers.CreateCommentOnProject)
	router.POST("/comments/for-comment/:comment_id", handlers.CreateCommentOnComment)
	router.GET("/comments/:comment_id", handlers.GetCommentById)
	router.PUT("/comments/:comment_id", handlers.UpdateCommentContent)
	router.DELETE("/comments/:comment_id", handlers.DeleteComment)

	router.GET("/comments/by-user/:user_id", handlers.GetCommentsByUserId)
	router.GET("/comments/by-post/:post_id", handlers.GetCommentsByPostId)
	router.GET("/comments/by-project/:project_id", handlers.GetCommentsByProjectId)
	router.GET("/comments/by-comment/:comment_id", handlers.GetCommentsByCommentId)

	router.POST("/comments/:username/likes/:comment_id", handlers.LikeComment)
	router.POST("/comments/:username/unlikes/:comment_id", handlers.UnlikeComment)
	router.GET("/comments/does-like/:username/:comment_id", handlers.IsCommentLiked)

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
