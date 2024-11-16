package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Some of the structs that will be used to interface between
// frontend, api, and db. Go has good handling for this

type User struct {
	Username     string    `json:"username"`
	Bio          string    `json:"bio"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
	// TODO: handle profile pic, followers, following
}

type Project struct {
	ID           int64     `json:"id"`
	Owner        int64     `json:"owner"` // linked with a User
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       int16     `json:"status"` // TODO: Enum? Does Go have Enums?
	Likes        int64     `json:"likes"`
	Tags         []string  `json:"tags"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
	// TODO: handle followers, posts
}

type Post struct {
	ID           int64     `json:"id"`
	User         int64     `json:"user"`    // Linked to User
	Project      int64     `json:"project"` // Linked to Project
	Likes        int64     `json:"likes"`
	Content      int64     `json:"content"`
	Comments     int64     `json:"comment"`
	CreationDate time.Time `json:"created_on"`
}

type Comment struct {
	ID            int64     `json:"id"`
	User          int64     `json:"user"`           // Linked to User
	Post          int64     `json:"post"`           // linked to Post
	ParentComment int64     `json:"parent_comment"` // TODO: does this need a Nil?
	CreationDate  time.Time `json:"created_on"`
}

// Test data, planning on creating maybe some sqlite
// so that devs can test better, not implemented rn

var testUsers = []User{
	{
		Username:     "dev_user1",
		Bio:          "Full-stack developer passionate about open-source projects.",
		Links:        []string{"https://github.com/dev_user1", "https://devuser1.com"},
		CreationDate: time.Now().AddDate(-1, 0, 0), // 1 year ago
	},
	{
		Username:     "tech_writer2",
		Bio:          "Technical writer and Python enthusiast.",
		Links:        []string{"https://blog.techwriter.com"},
		CreationDate: time.Now().AddDate(-2, 0, 0), // 2 years ago
	},
}

var testProjects = []Project{
	{
		ID:           1,
		Owner:        1,
		Name:         "OpenAPI Toolkit",
		Description:  "A toolkit for generating and testing OpenAPI specs.",
		Status:       1, // Assuming 1 = Active
		Likes:        120,
		Tags:         []string{"OpenAPI", "Go", "Tooling"},
		Links:        []string{"https://github.com/dev_user1/openapi-toolkit"},
		CreationDate: time.Now().AddDate(-1, -6, 0), // 1.5 years ago
	},
	{
		ID:           2,
		Owner:        2,
		Name:         "DocuHelper",
		Description:  "A library for streamlining technical documentation processes.",
		Status:       2, // Assuming 2 = Archived
		Likes:        85,
		Tags:         []string{"Documentation", "Python"},
		Links:        []string{"https://github.com/tech_writer2/docuhelper"},
		CreationDate: time.Now().AddDate(-3, 0, 0), // 3 years ago
	},
}

var testPosts = []Post{
	{
		ID:           1,
		User:         1,
		Project:      1,
		Likes:        40,
		Content:      1001, // Could reference a placeholder in a database
		Comments:     3,
		CreationDate: time.Now().AddDate(0, -3, 0), // 3 months ago
	},
	{
		ID:           2,
		User:         2,
		Project:      2,
		Likes:        25,
		Content:      1002,
		Comments:     2,
		CreationDate: time.Now().AddDate(0, -6, 0), // 6 months ago
	},
}

var testComments = []Comment{
	{
		ID:            1,
		User:          2,
		Post:          1,
		ParentComment: 0,                             // Root-level comment
		CreationDate:  time.Now().AddDate(0, -2, -5), // 2 months, 5 days ago
	},
	{
		ID:            2,
		User:          1,
		Post:          1,
		ParentComment: 1,                            // Reply to Comment 1
		CreationDate:  time.Now().AddDate(0, -1, 0), // 1 month ago
	},
}

// Here are all the routes, marked as v1 in case we need 
// backwards compatibility? Not sure, can change later

func getUsersV1(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, testUsers)
}

func createUserV1(context *gin.Context) {
	var newUser User
	err := context.BindJSON(&newUser)

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	testUsers = append(testUsers, newUser) // this would be the db
	context.IndentedJSON(http.StatusCreated, newUser)
}

func deleteUserV1(context *gin.Context) {
	username := context.Param("username")
	for i, entry := range testUsers {
		if entry.Username == username {
			// take the slice from up to the item we want to remove,
			// then append all of the items (unwound), from the
			// item after the one to be removed, not sure if this
			// is typical go pattern or not, but its pretty cool
			testUsers = append(testUsers[:i], testUsers[i+1:]...)
			context.IndentedJSON(http.StatusNoContent, gin.H{
				"message": fmt.Sprintf("Successfully deleted user %v.", username),
			}) // this doesnt return this message, at least in curl, is that an issue?
		}
	}
}

func updateUserInfoV1(context *gin.Context) {
	username := context.Param("username")

	var user *User
	// Find the user by username
	for i, entry := range testUsers {
		if entry.Username == username {
			user = &testUsers[i] // Use the address of the user in the slice
			break
		}
	}

	// If the user is not found, return 404
	if user == nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("User with username %v not found", username),
		})
		return
	}

	// Parse the incoming JSON data into a map
	var updateData map[string]interface{}
	err := context.BindJSON(&updateData)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	// Loop through the fields in the updateData map
	for key, value := range updateData {
		switch key {
		case "bio":
			// Ensure that bio is a string
			if bio, ok := value.(string); ok {
				user.Bio = bio
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid bio format"})
				return
			}
		case "links":
			// Ensure that links is an array of strings
			if links, ok := value.([]interface{}); ok {
				// Create a new string slice for links
				newLinks := []string{}
				for _, link := range links {
					// Check each link and assert that it is a string
					if strLink, ok := link.(string); ok {
						newLinks = append(newLinks, strLink)
					} else {
						context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid link format"})
						return
					}
				}
				user.Links = newLinks
			} else {
				context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid links format"})
				return
			}
		default:

			context.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid fields: %v", key),
			})
			return
		}
	}

	context.IndentedJSON(http.StatusOK, user)
}

func getUserByUsernameV1(context *gin.Context) {
	username := context.Param("username")

	for _, entry := range testUsers {
		if entry.Username == username {
			context.IndentedJSON(http.StatusOK, entry)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound,
		gin.H{"message": fmt.Sprintf("User with username %v not found", username)})
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsersV1)
	router.GET("/users/:username", getUserByUsernameV1)
	router.PUT("/users/:username", updateUserInfoV1)
	router.DELETE("/users/:username", deleteUserV1)
	router.POST("/users", createUserV1)

	router.Run("localhost:8000")
}
