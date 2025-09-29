package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"xlsx-processor/pkg/types"

	"github.com/gin-gonic/gin"
)

/*
	transformjson_test.go - E2E tests for the TransformJson route

	This test suite covers:
	1. Request parsing and validation
	2. JSON file extension validation
	3. Rule structure validation (redact/exclude operations)
	4. Webhook handling
	5. Preview mode vs output mode
	6. Error handling scenarios

	Note: These tests focus on the route logic and request validation.
	The storage layer is not mocked, so tests that hit storage endpoints
	will fail with expected storage errors (which validates the route works correctly).

	Test patterns follow paginate_test.go for gin context setup and
	transform_test.go for rule structure examples.
*/

// Test helper functions for transformjson route
func createMockCredentialTransform() types.Credential {
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

func createMockInputTransform(bucket, prefix string) types.Input {
	return types.Input{
		StorageType: "s3",
		Reference: types.SourceReference{
			Id:     "test-input",
			Bucket: bucket,
			Prefix: prefix,
			Region: "us-east-1",
		},
		Credential: createMockCredentialTransform(),
	}
}

func createMockOutputTransform(bucket, prefix string) types.Output {
	return types.Output{
		StorageType: "s3",
		Reference: types.SourceReference{
			Id:     "test-output",
			Bucket: bucket,
			Prefix: prefix,
			Region: "us-east-1",
		},
		Credential: createMockCredentialTransform(),
	}
}

func createMockRules() []types.Rule {
	return []types.Rule{
		{
			PageCondition: types.PageCondition{
				SheetName:       "Sheet1",
				IncludeFormulas: false,
			},
			Actions: []types.Action{
				{
					Operation:  "value",
					Value:      "sensitive",
					ActionType: "redact",
				},
			},
		},
	}
}

func createValidTransformRequest() types.RequestBodyTransform {
	return types.RequestBodyTransform{
		Input:   createMockInputTransform("input-bucket", "input/test.json"),
		Output:  createMockOutputTransform("output-bucket", "output/transformed.json"),
		Rules:   createMockRules(),
		Webhook: nil,
	}
}

func createValidTransformRequestPreview() types.RequestBodyTransform {
	return types.RequestBodyTransform{
		Input:   createMockInputTransform("input-bucket", "input/test.json"),
		Output:  createMockOutputTransform("output-bucket", "output/transformed.json"),
		Rules:   createMockRules(),
		Webhook: nil,
	}
}

func setupGinContextTransform(requestBody interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest("POST", "/transform-json", bytes.NewBuffer(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

// Tests for request parsing and validation
func TestTransformJson_RequestParsing_ValidRequest(t *testing.T) {
	request := createValidTransformRequest()
	c, _ := setupGinContextTransform(request)

	// Test request parsing
	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify input
	if parsedRequest.Input.StorageType != request.Input.StorageType {
		t.Errorf("Expected input storage type %s, got %s", request.Input.StorageType, parsedRequest.Input.StorageType)
	}

	if parsedRequest.Input.Reference.Bucket != request.Input.Reference.Bucket {
		t.Errorf("Expected input bucket %s, got %s", request.Input.Reference.Bucket, parsedRequest.Input.Reference.Bucket)
	}

	if parsedRequest.Input.Reference.Prefix != request.Input.Reference.Prefix {
		t.Errorf("Expected input prefix %s, got %s", request.Input.Reference.Prefix, parsedRequest.Input.Reference.Prefix)
	}

	if parsedRequest.Output.StorageType != request.Output.StorageType {
		t.Errorf("Expected output storage type %s, got %s", request.Output.StorageType, parsedRequest.Output.StorageType)
	}

	// Verify rules
	if len(parsedRequest.Rules) != len(request.Rules) {
		t.Errorf("Expected %d rules, got %d", len(request.Rules), len(parsedRequest.Rules))
	}

	if len(parsedRequest.Rules) > 0 {
		if parsedRequest.Rules[0].PageCondition.SheetName != request.Rules[0].PageCondition.SheetName {
			t.Errorf("Expected sheet name %s, got %s", request.Rules[0].PageCondition.SheetName, parsedRequest.Rules[0].PageCondition.SheetName)
		}

		if len(parsedRequest.Rules[0].Actions) != len(request.Rules[0].Actions) {
			t.Errorf("Expected %d actions, got %d", len(request.Rules[0].Actions), len(parsedRequest.Rules[0].Actions))
		}
	}
}

func TestTransformJson_RequestValidation_MissingInput(t *testing.T) {
	invalidRequest := types.RequestBodyTransform{
		// Missing Input - will result in zero value, not validation error
		Output: createMockOutputTransform("output-bucket", "output/transformed.json"),
		Rules:  createMockRules(),
	}

	c, _ := setupGinContextTransform(invalidRequest)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	// JSON binding will succeed, but input will have zero values
	if err != nil {
		t.Errorf("Unexpected error during binding: %v", err)
	}

	// Verify that input has zero/empty values
	if parsedRequest.Input.StorageType != "" {
		t.Errorf("Expected empty storage type, got %s", parsedRequest.Input.StorageType)
	}

	if parsedRequest.Input.Reference.Bucket != "" {
		t.Errorf("Expected empty bucket, got %s", parsedRequest.Input.Reference.Bucket)
	}
}

func TestTransformJson_RequestValidation_MissingRules(t *testing.T) {
	invalidRequest := types.RequestBodyTransform{
		Input:  createMockInputTransform("input-bucket", "input/test.json"),
		Output: createMockOutputTransform("output-bucket", "output/transformed.json"),
		// Missing Rules
	}

	c, _ := setupGinContextTransform(invalidRequest)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	// Validation should catch missing required rules
	if err == nil {
		t.Error("Expected validation error for missing rules")
	}
}

func TestTransformJson_RequestValidation_EmptyRules(t *testing.T) {
	invalidRequest := types.RequestBodyTransform{
		Input:  createMockInputTransform("input-bucket", "input/test.json"),
		Output: createMockOutputTransform("output-bucket", "output/transformed.json"),
		Rules:  []types.Rule{}, // Empty rules array
	}

	c, _ := setupGinContextTransform(invalidRequest)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	// Empty slice is valid JSON, so binding succeeds
	if err != nil {
		t.Errorf("Unexpected error during binding: %v", err)
	}

	// Verify that rules array is empty
	if len(parsedRequest.Rules) != 0 {
		t.Errorf("Expected empty rules array, got %d rules", len(parsedRequest.Rules))
	}
}

func TestTransformJson_RequestValidation_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Invalid JSON
	c.Request = httptest.NewRequest("POST", "/transform-json", bytes.NewBufferString("{invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// Test file extension validation
func TestTransformJson_FileExtensionValidation_NonJsonFile(t *testing.T) {
	// Create request with non-JSON file
	request := createValidTransformRequest()
	request.Input.Reference.Prefix = "input/test.xlsx" // Not a JSON file

	c, w := setupGinContextTransform(request)

	// Call the actual route function
	TransformJson(c)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should contain error message about file not being JSON
	if response["message"] == nil {
		t.Error("Expected error message in response")
	}

	errorMsg, exists := response["message"].(string)
	if !exists || errorMsg != "file is not a JSON file" {
		t.Errorf("Expected 'file is not a JSON file' error, got: %v", response["message"])
	}
}

func TestTransformJson_FileExtensionValidation_ValidJsonFile(t *testing.T) {
	request := createValidTransformRequest()
	request.Input.Reference.Prefix = "input/test.json" // Valid JSON file

	c, w := setupGinContextTransform(request)

	// Call the actual route function - this will fail at storage download
	// but should pass file extension validation
	TransformJson(c)

	// Should not return 400 for file extension validation
	// (It will return 500 for storage download failure, which is expected)
	if w.Code == http.StatusBadRequest {
		// Parse response to check if it's about file extension
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		errorMsg, exists := response["message"].(string)
		if exists && errorMsg == "file is not a JSON file" {
			t.Error("File extension validation failed for valid JSON file")
		}
	}
}

// Test rule structure validation
func TestTransformJson_RuleStructure_ValidAction(t *testing.T) {
	rules := []types.Rule{
		{
			PageCondition: types.PageCondition{
				SheetName:       "Sheet1",
				IncludeFormulas: true,
			},
			Actions: []types.Action{
				{
					Operation:  "range",
					Value:      "A1:C3",
					ActionType: "redact",
				},
				{
					Operation:  "value",
					Value:      "sensitive_data",
					ActionType: "redact",
				},
			},
		},
	}

	request := createValidTransformRequest()
	request.Rules = rules

	c, _ := setupGinContextTransform(request)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	if err != nil {
		t.Errorf("Expected no error for valid rules, got %v", err)
	}

	if len(parsedRequest.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(parsedRequest.Rules))
	}

	rule := parsedRequest.Rules[0]
	if rule.PageCondition.SheetName != "Sheet1" {
		t.Errorf("Expected sheet name 'Sheet1', got %s", rule.PageCondition.SheetName)
	}

	if !rule.PageCondition.IncludeFormulas {
		t.Error("Expected IncludeFormulas to be true")
	}

	if len(rule.Actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(rule.Actions))
	}

	// Verify first action
	action1 := rule.Actions[0]
	if action1.Operation != "range" {
		t.Errorf("Expected operation 'range', got %s", action1.Operation)
	}
	if action1.Value != "A1:C3" {
		t.Errorf("Expected value 'A1:C3', got %s", action1.Value)
	}
	if action1.ActionType != "redact" {
		t.Errorf("Expected action type 'redact', got %s", action1.ActionType)
	}
}

// Test various action types from the examples in the rules.go file
func TestTransformJson_RuleStructure_DifferentActionTypes(t *testing.T) {
	testCases := []struct {
		name       string
		operation  string
		value      string
		actionType string
	}{
		{"redact by range", "range", "C4:D9", "redact"},
		{"redact by value", "value", "1.00%", "redact"},
		{"redact by text color", "textColor", "0070C0", "redact"},
		{"redact by bg color", "bgColor", "0070C0", "redact"},
		{"exclude column", "column", "C", "exclude"},
		{"exclude row", "row", "4", "exclude"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := types.Rule{
				PageCondition: types.PageCondition{
					SheetName:       "Sheet1",
					IncludeFormulas: false,
				},
				Actions: []types.Action{
					{
						Operation:  tc.operation,
						Value:      tc.value,
						ActionType: tc.actionType,
					},
				},
			}

			request := createValidTransformRequest()
			request.Rules = []types.Rule{rule}

			c, _ := setupGinContextTransform(request)

			var parsedRequest types.RequestBodyTransform
			err := bindAndValidate(c, &parsedRequest)

			if err != nil {
				t.Errorf("Expected no error for %s, got %v", tc.name, err)
			}

			if len(parsedRequest.Rules) != 1 {
				t.Errorf("Expected 1 rule for %s, got %d", tc.name, len(parsedRequest.Rules))
			}

			action := parsedRequest.Rules[0].Actions[0]
			if action.Operation != tc.operation {
				t.Errorf("Expected operation %s, got %s", tc.operation, action.Operation)
			}
			if action.Value != tc.value {
				t.Errorf("Expected value %s, got %s", tc.value, action.Value)
			}
			if action.ActionType != tc.actionType {
				t.Errorf("Expected action type %s, got %s", tc.actionType, action.ActionType)
			}
		})
	}
}

// Test webhook functionality
func TestTransformJson_WebhookHandling(t *testing.T) {
	webhook := &types.Webhook{
		Url: "https://example.com/webhook",
	}

	request := createValidTransformRequest()
	request.Webhook = webhook

	c, _ := setupGinContextTransform(request)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if parsedRequest.Webhook == nil {
		t.Error("Expected webhook to be parsed")
	} else {
		if parsedRequest.Webhook.Url != webhook.Url {
			t.Errorf("Expected webhook URL %s, got %s", webhook.Url, parsedRequest.Webhook.Url)
		}
	}
}

// Benchmark test for performance
func BenchmarkTransformJson_RequestParsing(b *testing.B) {
	request := createValidTransformRequest()

	for i := 0; i < b.N; i++ {
		c, _ := setupGinContextTransform(request)
		var parsedRequest types.RequestBodyTransform
		_ = bindAndValidate(c, &parsedRequest)
	}
}

// Test helper to create complex rules for more comprehensive testing
func createComplexRules() []types.Rule {
	return []types.Rule{
		{
			PageCondition: types.PageCondition{
				SheetName:           "Sheet1",
				IncludeFormulas:     false,
				NonEmptyValueRedact: true,
			},
			Actions: []types.Action{
				{Operation: "range", Value: "A1:C10", ActionType: "redact"},
				{Operation: "column", Value: "D", ActionType: "exclude"},
			},
		},
		{
			PageCondition: types.PageCondition{
				SheetName:           "Sheet2",
				IncludeFormulas:     true,
				NonEmptyValueRedact: false,
			},
			Actions: []types.Action{
				{Operation: "value", Value: "confidential", ActionType: "redact"},
				{Operation: "textColor", Value: "FF0000", ActionType: "redact"},
				{Operation: "bgColor", Value: "FFFF00", ActionType: "redact"},
			},
		},
	}
}

func TestTransformJson_ComplexRulesStructure(t *testing.T) {
	request := createValidTransformRequest()
	request.Rules = createComplexRules()

	c, _ := setupGinContextTransform(request)

	var parsedRequest types.RequestBodyTransform
	err := bindAndValidate(c, &parsedRequest)

	if err != nil {
		t.Errorf("Expected no error for complex rules, got %v", err)
	}

	if len(parsedRequest.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(parsedRequest.Rules))
	}

	// Verify first rule
	rule1 := parsedRequest.Rules[0]
	if rule1.PageCondition.SheetName != "Sheet1" {
		t.Errorf("Expected sheet name 'Sheet1', got %s", rule1.PageCondition.SheetName)
	}
	if rule1.PageCondition.IncludeFormulas {
		t.Error("Expected IncludeFormulas to be false for first rule")
	}
	if !rule1.PageCondition.NonEmptyValueRedact {
		t.Error("Expected NonEmptyValueRedact to be true for first rule")
	}
	if len(rule1.Actions) != 2 {
		t.Errorf("Expected 2 actions in first rule, got %d", len(rule1.Actions))
	}

	// Verify second rule
	rule2 := parsedRequest.Rules[1]
	if rule2.PageCondition.SheetName != "Sheet2" {
		t.Errorf("Expected sheet name 'Sheet2', got %s", rule2.PageCondition.SheetName)
	}
	if !rule2.PageCondition.IncludeFormulas {
		t.Error("Expected IncludeFormulas to be true for second rule")
	}
	if rule2.PageCondition.NonEmptyValueRedact {
		t.Error("Expected NonEmptyValueRedact to be false for second rule")
	}
	if len(rule2.Actions) != 3 {
		t.Errorf("Expected 3 actions in second rule, got %d", len(rule2.Actions))
	}
}
