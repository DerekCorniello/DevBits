package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Method         string
	Endpoint       string
	Input          string
	ExpectedStatus int
	ExpectedBody   string
}

var main_tests []TestCase = []TestCase{
	{
		Method:         http.MethodGet,
		Endpoint:       "/health",
		Input:          "",
		ExpectedStatus: http.StatusOK,
		ExpectedBody:   `{"message":"API is running!"}`,
	},
	{
		Method:         http.MethodPost,
		Endpoint:       "/health",
		Input:          `{"test": "data"}`,
		ExpectedStatus: http.StatusMethodNotAllowed,
		ExpectedBody:   `405 method not allowed`,
	},
}

func (tc *TestCase) Run(t *testing.T) {
	t.Helper() // Mark this function as a test helper

	url := "http://localhost:8080" + tc.Endpoint
	var req *http.Request
	var err error

	// set the request body if it is provided
	if tc.Input != "" {
		req, err = http.NewRequest(tc.Method, url, bytes.NewBufferString(tc.Input))
	} else {
		req, err = http.NewRequest(tc.Method, url, nil)
	}
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	assert.Equal(t, tc.ExpectedStatus, resp.StatusCode, "Status code mismatch for %s %s", tc.Method, tc.Endpoint)

	// check if the response is expected to be JSON
	if resp.Header.Get("Content-Type") == "application/json" {
		// try to parse response JSON
		var actualJSON map[string]interface{}
		if err := json.Unmarshal(body, &actualJSON); err != nil {
			t.Fatalf("Expected valid JSON response but got invalid JSON. Body: %q, Error: %v", body, err)
		}

		// parse the expected JSON
		var expectedJSON map[string]interface{}
		if err := json.Unmarshal([]byte(tc.ExpectedBody), &expectedJSON); err != nil {
			t.Fatalf("Test has invalid ExpectedBody JSON: %q, Error: %v", tc.ExpectedBody, err)
		}

		assert.Equal(t, expectedJSON, actualJSON, "Response body mismatch for %s %s", tc.Method, tc.Endpoint)
	} else {
		assert.Equal(t, tc.ExpectedBody, string(body), "Response body mismatch for %s %s", tc.Method, tc.Endpoint)
	}
}

func TestAPI(t *testing.T) {
	// put all tests cases together
	tests := append(main_tests, append(user_tests, append(project_tests, append(comment_tests, post_tests...)...)...)...)

	for _, test := range tests {
		t.Run(test.Method+" "+test.Endpoint, func(t *testing.T) {
			test.Run(t)
		})
	}
}
