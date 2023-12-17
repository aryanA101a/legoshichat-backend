package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	// Create a new instance of Error
	err := NewError("Test error message")

	// Check for expectations
	assert.Equal(t, "Test error message", err.Error(), "Mismatch in error message")
}

func TestHttpStatusError(t *testing.T) {
	// Test cases for different HTTP status codes
	testCases := []struct {
		statusCode    int
		expectedError string
	}{
		{400, "Bad Request - Test message"},
		{401, "Unauthorized - Test message"},
		{403, "Forbidden - Test message"},
		{404, "Not Found - Test message"},
		{500, "Internal Server Error - Test message"},
		{999, "Unknown Status Code - Test message"},
	}

	// Iterate through test cases
	for _, testCase := range testCases {
		// Call HttpStatusError with the test case parameters
		httpError := HttpStatusError(testCase.statusCode, "Test message")

		// Check for expectations
		assert.Equal(t, testCase.statusCode, httpError.StatusCode, "Mismatch in HTTP status code")
		assert.Equal(t, testCase.expectedError, httpError.Err.Error(), "Mismatch in HTTP error message")
	}
}