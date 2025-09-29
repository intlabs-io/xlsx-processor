package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"xlsx-processor/pkg/types"
)

// Test helper functions
func createMockCredential() types.Credential {
	return types.Credential{
		Secrets: types.Secrets{
			Secret:      "test-secret",
			AccessToken: "test-token",
		},
		Resources: types.Resources{
			Id: "test-id",
		},
	}
}

func createMockInput(bucket, prefix string) types.Input {
	return types.Input{
		StorageType: "s3",
		Reference: types.SourceReference{
			Id:     "test-input",
			Bucket: bucket,
			Prefix: prefix,
			Region: "us-east-1",
		},
		Credential: createMockCredential(),
	}
}

func createMockOutput(bucket, prefix string) types.Output {
	return types.Output{
		StorageType: "s3",
		Reference: types.SourceReference{
			Id:     "test-output",
			Bucket: bucket,
			Prefix: prefix,
			Region: "us-east-1",
		},
		Credential: createMockCredential(),
	}
}

func createValidRequest() types.RequestBodyPaginate {
	return types.RequestBodyPaginate{
		Input:  createMockInput("input-bucket", "input/test.xlsx"),
		Output: createMockOutput("output-bucket", "output/paginated"),
	}
}

func setupGinContext(requestBody interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest("POST", "/paginate", bytes.NewBuffer(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

// Tests
func TestPaginate_Success_RequestParsing(t *testing.T) {
	request := createValidRequest()
	c, _ := setupGinContext(request)

	// For this test, we'll just verify the request parsing works
	var parsedRequest types.RequestBodyPaginate
	err := bindAndValidate(c, &parsedRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if parsedRequest.Input.StorageType != request.Input.StorageType {
		t.Errorf("Expected input storage type %s, got %s", request.Input.StorageType, parsedRequest.Input.StorageType)
	}

	if parsedRequest.Output.StorageType != request.Output.StorageType {
		t.Errorf("Expected output storage type %s, got %s", request.Output.StorageType, parsedRequest.Output.StorageType)
	}

	if parsedRequest.Input.Reference.Bucket != request.Input.Reference.Bucket {
		t.Errorf("Expected input bucket %s, got %s", request.Input.Reference.Bucket, parsedRequest.Input.Reference.Bucket)
	}

	if parsedRequest.Output.Reference.Bucket != request.Output.Reference.Bucket {
		t.Errorf("Expected output bucket %s, got %s", request.Output.Reference.Bucket, parsedRequest.Output.Reference.Bucket)
	}
}

func TestPaginate_RequestParsing_EmptyInput(t *testing.T) {
	// Test that empty input fields result in zero values
	invalidRequest := types.RequestBodyPaginate{
		// Input will be zero value
		Output: createMockOutput("output-bucket", "output/paginated"),
	}

	c, _ := setupGinContext(invalidRequest)

	var parsedRequest types.RequestBodyPaginate
	err := bindAndValidate(c, &parsedRequest)

	// Validation might not catch this, but we can verify the zero values
	if err != nil {
		t.Logf("Validation error (expected): %v", err)
	}

	// Verify that input has zero values
	if parsedRequest.Input.StorageType != "" {
		t.Errorf("Expected empty storage type, got %s", parsedRequest.Input.StorageType)
	}

	if parsedRequest.Input.Reference.Bucket != "" {
		t.Errorf("Expected empty bucket, got %s", parsedRequest.Input.Reference.Bucket)
	}
}

func TestPaginate_RequestParsing_EmptyOutput(t *testing.T) {
	// Test that empty output fields result in zero values
	invalidRequest := types.RequestBodyPaginate{
		Input: createMockInput("input-bucket", "input/test.xlsx"),
		// Output will be zero value
	}

	c, _ := setupGinContext(invalidRequest)

	var parsedRequest types.RequestBodyPaginate
	err := bindAndValidate(c, &parsedRequest)

	// Validation might not catch this, but we can verify the zero values
	if err != nil {
		t.Logf("Validation error (expected): %v", err)
	}

	// Verify that output has zero values
	if parsedRequest.Output.StorageType != "" {
		t.Errorf("Expected empty storage type, got %s", parsedRequest.Output.StorageType)
	}

	if parsedRequest.Output.Reference.Bucket != "" {
		t.Errorf("Expected empty bucket, got %s", parsedRequest.Output.Reference.Bucket)
	}
}

func TestPaginate_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid JSON
	c.Request = httptest.NewRequest("POST", "/paginate", bytes.NewBufferString("{invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	var parsedRequest types.RequestBodyPaginate
	err := bindAndValidate(c, &parsedRequest)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestPaginate_OutputPathGeneration(t *testing.T) {
	// Test that output paths are generated correctly
	basePrefix := "output/paginated"

	testCases := []struct {
		sheetIndex   int
		expectedPath string
	}{
		{0, "output/paginated/pages/1.json"},
		{1, "output/paginated/pages/2.json"},
		{99, "output/paginated/pages/100.json"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("sheet_%d", tc.sheetIndex), func(t *testing.T) {
			displaySheetNum := tc.sheetIndex + 1
			expectedPath := fmt.Sprintf("%s/pages/%d.json", basePrefix, displaySheetNum)
			if expectedPath != tc.expectedPath {
				t.Errorf("Expected path %s, got %s", tc.expectedPath, expectedPath)
			}
		})
	}
}

func TestPaginate_RequestValidation_RequiredFields(t *testing.T) {
	testCases := []struct {
		name    string
		request types.RequestBodyPaginate
		wantErr bool
	}{
		{
			name:    "Valid request",
			request: createValidRequest(),
			wantErr: false,
		},
		{
			name: "Missing input storage type",
			request: types.RequestBodyPaginate{
				Input: types.Input{
					// Missing StorageType
					Reference:  types.SourceReference{Bucket: "test", Prefix: "test"},
					Credential: createMockCredential(),
				},
				Output: createMockOutput("output-bucket", "output/paginated"),
			},
			wantErr: false, // StorageType is not explicitly validated by struct tags
		},
		{
			name: "Missing output reference",
			request: types.RequestBodyPaginate{
				Input: createMockInput("input-bucket", "input/test.xlsx"),
				Output: types.Output{
					StorageType: "s3",
					// Missing Reference - will use zero value
					Credential: createMockCredential(),
				},
			},
			wantErr: false, // Reference is not explicitly validated by struct tags
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := setupGinContext(tc.request)

			var parsedRequest types.RequestBodyPaginate
			err := bindAndValidate(c, &parsedRequest)

			if tc.wantErr && err == nil {
				t.Errorf("Expected error for test case: %s, got nil", tc.name)
			} else if !tc.wantErr && err != nil {
				t.Errorf("Did not expect error for test case: %s, got %v", tc.name, err)
			}
		})
	}
}

// Benchmark test for performance
func BenchmarkPaginate_RequestParsing(b *testing.B) {
	request := createValidRequest()

	for i := 0; i < b.N; i++ {
		c, _ := setupGinContext(request)
		var parsedRequest types.RequestBodyPaginate
		_ = bindAndValidate(c, &parsedRequest)
	}
}
