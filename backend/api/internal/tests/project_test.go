package tests

import (
	"net/http"
)

var project_tests []TestCase = []TestCase{
	// Test GET by project ID
	{
		Method:         http.MethodGet,
		Endpoint:       "/projects/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"id":1,"owner":1,"name":"OpenAPI Toolkit","description":"A toolkit for generating and testing OpenAPI specs.","status":1,"likes":120,"tags":["OpenAPI","Go","Tooling"],"links":["https://github.com/dev_user1/openapi-toolkit"],"creation_date":"2023-06-13T00:00:00Z"}`,
	},
	{
		Method:         http.MethodGet,
		Endpoint:       "/projects/-1",
		Input:          "",
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Project with id '-1' not found"}`,
	},

	// Test POST to create a new project
	{
		Method:         http.MethodPost,
		Endpoint:       "/projects",
		Input:          `{"name":"New Project","description":"Test project description","owner":1}`,
		ExpectedStatus: http.StatusCreated,
		ExpectedBody:   `{"message":"Project created successfully with id '4'"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/projects",
		Input:          `{"name":"","description":"Test project description","owner":1}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Failed to bind to JSON: Key: 'Project.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/projects",
		Input:          `{"name":"Duplicate Project","description":"Test duplicate","owner":-1}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Failed to verify project ownership. User could not be found"}`,
	},

	// Test PUT to update project info
	{
		Method:         http.MethodPut,
		Endpoint:       "/projects/4",
		Input:          `{"name":"Updated Project","description":"Updated description"}`,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"Project updated successfully","project":4}`,
	},
	{
		Method:         http.MethodPut,
		Endpoint:       "/projects/4",
		Input:          `{"owner":9999}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Invalid owner id: 9999"}`,
	},
	{
		Method:         http.MethodPut,
		Endpoint:       "/projects/9999",
		Input:          `{"name":"Non-existent Project"}`,
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Project with id '9999' not found"}`,
	},

	// Test DELETE project
	{
		Method:         http.MethodDelete,
		Endpoint:       "/projects/4",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"Project 4 deleted."}`,
	},
	{
		Method:         http.MethodDelete,
		Endpoint:       "/projects/9999",
		Input:          "",
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Failed to delete project: Deletion did not affect any records"}`,
	},
}
