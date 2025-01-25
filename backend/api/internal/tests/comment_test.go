package tests

import (
	"net/http"
)

var comment_tests []TestCase = []TestCase{
	// Test if comment is editable
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/can-edit/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"status":false}`,
	},
	// Test GET comment by ID
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"id":1,"user":1,"likes":5,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"This is a fantastic project! Can't wait to contribute."}`,
	},
	// Test GET non-existent comment
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/-9999",
		Input:          "",
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Comment with id -9999 not found"}`,
	},
	// Test CREATE comment on post
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/for-post/1",
		Input:          `{"user":1,"content":"New comment on post","parent_comment":null}`,
		ExpectedStatus: http.StatusCreated,
		ExpectedBody:   `{"message":"Comment created successfully with id 13"}`,
	},
	// Test CREATE comment on project
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/for-project/1",
		Input:          `{"user":2,"content":"New comment on project","parent_comment":null}`,
		ExpectedStatus: http.StatusCreated,
		ExpectedBody:   `{"message":"Comment created successfully with id 14"}`,
	},
	// Test CREATE reply to comment
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/for-comment/1",
		Input:          `{"user":3,"content":"Reply to existing comment","parent_comment":1}`,
		ExpectedStatus: http.StatusCreated,
		ExpectedBody:   `{"message":"Reply created successfully with id 15"}`,
	},
	// Test if comment is editable
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/can-edit/15",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"status":true}`,
	},
	// Test UPDATE comment
	{
		Method:         http.MethodPut,
		Endpoint:       "/comments/15",
		Input:          `{"content":"Updated comment content"}`,
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"comment":{"content":"Updated comment content","id":15,"likes":0,"parent_comment":1,"user":3},"message":"Comment updated successfully"}`,
	},
	// Test bad UPDATE comment
	{
		Method:         http.MethodPut,
		Endpoint:       "/comments/1",
		Input:          `{"content":"Updated comment content"}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Error updating comment: Cannot update comment. More than 2 minutes have passed since posting."}`,
	},
	// Test DELETE comment
	{
		Method:         http.MethodDelete,
		Endpoint:       "/comments/13",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"Comment 13 soft deleted."}`,
	},
	{
		Method:         http.MethodDelete,
		Endpoint:       "/comments/14",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"Comment 14 soft deleted."}`,
	},
	{
		Method:         http.MethodDelete,
		Endpoint:       "/comments/15",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"Comment 15 soft deleted."}`,
	},
	// Test GET comments by user
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/by-user/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[{"id":1,"user":1,"likes":1,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"This is a fantastic project! Can't wait to contribute."},{"id":2,"user":2,"likes":0,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"I love the concept, but I think the documentation could be improved."},{"id":3,"user":4,"likes":1,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Great to see more open-source tools for API development!"},{"id":4,"user":3,"likes":1,"parent_comment":3,"created_on":"2024-12-23T00:00:00Z","content":"I agree, but the API specs seem a bit too complex for beginners."},{"id":5,"user":5,"likes":0,"parent_comment":1,"created_on":"2024-12-23T00:00:00Z","content":"I hope this toolkit will integrate with other Go tools soon!"},{"id":6,"user":3,"likes":0,"parent_comment":2,"created_on":"2024-12-23T00:00:00Z","content":"I agree, the documentation is lacking in detail."},{"id":14,"user":-1,"likes":0,"parent_comment":null,"created_on":"1970-01-01T00:00:00Z","content":"This comment was deleted."},{"id":12,"user":1,"likes":2,"parent_comment":3,"created_on":"2024-12-23T00:00:00Z","content":"Looking forward to testing it!"}]`,
	},
	// Test GET comments by post
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/by-post/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[{"id":7,"user":4,"likes":2,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Awesome update! I'll try it out."},{"id":8,"user":3,"likes":1,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Thanks for sharing! Will this feature be extended soon?"},{"id":9,"user":5,"likes":4,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Great work, looking forward to more updates!"},{"id":10,"user":2,"likes":1,"parent_comment":2,"created_on":"2024-12-23T00:00:00Z","content":"Will this be compatible with earlier versions of OpenAPI?"},{"id":11,"user":3,"likes":3,"parent_comment":1,"created_on":"2024-12-23T00:00:00Z","content":"I hope the next update addresses performance improvements."},{"id":12,"user":1,"likes":2,"parent_comment":3,"created_on":"2024-12-23T00:00:00Z","content":"Looking forward to testing it!"},{"id":13,"user":-1,"likes":0,"parent_comment":null,"created_on":"1970-01-01T00:00:00Z","content":"This comment was deleted."}]`,
	},
	// Test GET comments by project
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/by-project/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[{"id":1,"user":1,"likes":5,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"This is a fantastic project! Can't wait to contribute."},{"id":2,"user":2,"likes":3,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"I love the concept, but I think the documentation could be improved."},{"id":3,"user":4,"likes":4,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Great to see more open-source tools for API development!"},{"id":4,"user":3,"likes":2,"parent_comment":3,"created_on":"2024-12-23T00:00:00Z","content":"I agree, but the API specs seem a bit too complex for beginners."},{"id":5,"user":5,"likes":1,"parent_comment":1,"created_on":"2024-12-23T00:00:00Z","content":"I hope this toolkit will integrate with other Go tools soon!"},{"id":6,"user":3,"likes":1,"parent_comment":2,"created_on":"2024-12-23T00:00:00Z","content":"I agree, the documentation is lacking in detail."},{"id":14,"user":-1,"likes":0,"parent_comment":null,"created_on":"1970-01-01T00:00:00Z","content":"This comment was deleted."}]`,
	},
	// Test GET replies to comment
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/by-comment/3",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `[{"id":8,"user":3,"likes":1,"parent_comment":null,"created_on":"2024-12-23T00:00:00Z","content":"Thanks for sharing! Will this feature be extended soon?"},{"id":11,"user":3,"likes":3,"parent_comment":1,"created_on":"2024-12-23T00:00:00Z","content":"I hope the next update addresses performance improvements."}]`,
	},
	// Test LIKE comment
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/dev_user1/likes/2",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"dev_user1 likes comment 2"}`,
	},
	// Test UNLIKE comment
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/dev_user1/unlikes/2",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"dev_user1 unliked comment 2"}`,
	},
	// Test CHECK if comment is liked
	{
		Method:         http.MethodGet,
		Endpoint:       "/comments/does-like/tech_writer2/1",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"status":true}`,
	},
	// Test invalid inputs
	{
		Method:         http.MethodPost,
		Endpoint:       "/comments/for-post/1",
		Input:          `{"invalid":"json"}`,
		ExpectedStatus: http.StatusBadRequest,
		ExpectedBody:   `{"error":"Bad Request","message":"Failed to bind to JSON: Key: 'Comment.User' Error:Field validation for 'User' failed on the 'required' tag\nKey: 'Comment.Content' Error:Field validation for 'Content' failed on the 'required' tag"}`,
	},
	{
		Method:         http.MethodPut,
		Endpoint:       "/comments/-9999",
		Input:          `{"content":"Update non-existent comment"}`,
		ExpectedStatus: http.StatusNotFound,
		ExpectedBody:   `{"error":"Not Found","message":"Comment with id -9999 not found"}`,
	},
}
