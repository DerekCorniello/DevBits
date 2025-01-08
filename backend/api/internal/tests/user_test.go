package tests

import (
	"net/http"
)

var user_tests []TestCase = []TestCase{

	// test get by user
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"username":"dev_user1","bio":"Full-stack developer passionate about open-source projects.","links":["https://github.com/dev_user1","https://devuser1.com"],"created_on":"2023-12-13T00:00:00Z","picture":"https://example.com/dev_user1.jpg"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/dev_user1",
		Input:          "",
		ExpectedStatus: http.StatusMethodNotAllowed,
		ExpectedBody:   `405 method not allowed`,
	},

	// test get by non-user
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/non_dev_user1",
		Input:          "",
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"User with username 'non_dev_user1' not found"}`,
	},

	// create user
	{
		Method:         http.MethodPost,
		Endpoint:       "/users",
		Input:          `{"username": "new_user","bio":"This is a test user.","links":["https://example.com","https://another-link.com"],"picture":"https://example.com/profile.jpg"}`,
		ExpectedStatus: http.StatusCreated,
		ExpectedBody:   `{"message":"Created new user: 'new_user'"}`,
	},

	// creating user with same name
	{
		Method:         http.MethodPost,
		Endpoint:       "/users",
		Input:          `{"username": "new_user","bio":"This is a test user.","links":["https://example.com","https://another-link.com"],"picture":"https://example.com/profile.jpg"}`,
		ExpectedStatus: http.StatusInternalServerError,
		ExpectedBody:   `{"error":"Internal Server Error","message":"Failed to create user: Failed to create user 'new_user': UNIQUE constraint failed: Users.username"}`,
	},

	// creating user with no name
	{
		Method:         http.MethodPost,
		Endpoint:       "/users",
		Input:          `{"username": "","bio":"This is a test user.","links":["https://example.com","https://another-link.com"],"picture":"https://example.com/profile.jpg"}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Failed to bind to JSON: Key: 'User.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
	},

	// updating user with no name
	{
		Method:         http.MethodPut,
		Endpoint:       "/users/new_user",
		Input:          `{"username": ""}`,
		ExpectedStatus: http.StatusInternalServerError,
		ExpectedBody:   `{"error":"Internal Server Error","message":"Error updating user: Updated username cannot be empty!"}`,
	},

	// updated user with name in use
	{
		Method:         http.MethodPut,
		Endpoint:       "/users/dev_user1",
		Input:          `{"username": "tech_writer2","bio":"This is a test user.","links":["https://example.com","https://another-link.com"],"picture":"https://example.com/profile.jpg"}`,
		ExpectedStatus: http.StatusInternalServerError,
		ExpectedBody:   `{"error":"Internal Server Error","message":"Error updating user: Error checking rows affected: Error executing update query: UNIQUE constraint failed: Users.username"}`,
	},

	// updating all user fields
	{
		Method:         http.MethodPut,
		Endpoint:       "/users/dev_user1",
		Input:          `{"username": "new_user_updated","bio":"This is the test user's updated bio.","links":["https://example.com/updated","https://another-link-updated.com"],"picture":"https://example.com/updates_profile.jpg"}`,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"User updated successfully.","user":{"username":"new_user_updated","bio":"This is the test user's updated bio.","links":["https://example.com/updated","https://another-link-updated.com"],"created_on":"2023-12-13T00:00:00Z","picture":"https://example.com/updates_profile.jpg"}}`,
	},

    // update it back...

	{
		Method:         http.MethodPut,
		Endpoint:       "/users/new_user_updated",
		Input:          `{"username":"dev_user1","bio":"Full-stack developer passionate about open-source projects.","links":["https://github.com/dev_user1","https://devuser1.com"],"created_on":"2023-12-13T00:00:00Z","picture":"https://example.com/dev_user1.jpg"}`,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"User updated successfully.","user":{"username":"dev_user1","bio":"Full-stack developer passionate about open-source projects.","links":["https://github.com/dev_user1","https://devuser1.com"],"created_on":"2023-12-13T00:00:00Z","picture":"https://example.com/dev_user1.jpg"}}`,
	},

	// delete our test user
	{
		Method:         http.MethodDelete,
		Endpoint:       "/users/new_user",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"User 'new_user' deleted."}`,
	},

	// delete non-existant user
	{
		Method:         http.MethodDelete,
		Endpoint:       "/users/not_a_user",
		Input:          ``,
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Failed to delete user: Deletion did not affect any records"}`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/follows/usernames",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"following":["data_scientist3"],"message":"Successfully got following"}`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/follows",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[3]`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/followers",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[4]`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/dev_user1/follow/tech_writer2",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
        ExpectedBody:   `{"message":"dev_user1 now follows tech_writer2"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/tech_writer2/follow/dev",
		Input:          ``,
		ExpectedStatus: http.StatusNotFound,
        ExpectedBody:   `{"error":"Not Found","message":"Failed to add follower: Cannot find user with username 'dev'"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/tech/follow/dev_user1",
		Input:          ``,
		ExpectedStatus: http.StatusNotFound,
        ExpectedBody:   `{"error":"Not Found","message":"Failed to add follower: Cannot find user with username 'tech'"}`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/follows",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[2,3]`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/followers",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[4]`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/dev_user1/unfollow/tech_writer2",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
        ExpectedBody:   `{"message":"dev_user1 unfollowed tech_writer2"}`,
	},	
    {
		Method:         http.MethodPost,
		Endpoint:       "/users/backend_guru4/unfollow/dev_user1",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
        ExpectedBody:   `{"message":"backend_guru4 unfollowed dev_user1"}`,
	},
    {
		Method:         http.MethodPost,
		Endpoint:       "/users/tech_writer2/unfollow/dev",
		Input:          ``,
		ExpectedStatus: http.StatusNotFound,
        ExpectedBody:   `{"error":"Not Found","message":"Failed to remove follower: Cannot find user with username 'dev'"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/users/tech/unfollow/dev_user1",
		Input:          ``,
		ExpectedStatus: http.StatusNotFound,
        ExpectedBody:   `{"error":"Not Found","message":"Failed to remove follower: Cannot find user with username 'tech'"}`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/follows",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[3]`,
	},
    {
		Method:         http.MethodGet,
		Endpoint:       "/users/dev_user1/followers",
		Input:          ``,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `null`,
	},
}
