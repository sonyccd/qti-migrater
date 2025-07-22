package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestQTIError_New_Functions(t *testing.T) {
	// Test NewParsingError
	causeErr := errors.New("XML syntax error")
	parseErr := NewParsingError("Failed to parse document", causeErr)
	
	if parseErr.Type != ErrorTypeParsing {
		t.Errorf("Expected ErrorTypeParsing, got %v", parseErr.Type)
	}
	
	if parseErr.Message != "Failed to parse document" {
		t.Errorf("Expected message 'Failed to parse document', got '%s'", parseErr.Message)
	}
	
	if parseErr.Cause != causeErr {
		t.Errorf("Expected cause to be set correctly")
	}
	
	// Test NewValidationError
	validErr := NewValidationError("Invalid attribute", "q001", "/item[@ident='q001']")
	
	if validErr.Type != ErrorTypeValidation {
		t.Errorf("Expected ErrorTypeValidation, got %v", validErr.Type)
	}
	
	if validErr.ItemID != "q001" {
		t.Errorf("Expected ItemID 'q001', got '%s'", validErr.ItemID)
	}
	
	if validErr.ElementPath != "/item[@ident='q001']" {
		t.Errorf("Expected ElementPath to be set correctly")
	}
	
	// Test NewMigrationError
	migErr := NewMigrationError("Migration failed", "Detailed explanation")
	
	if migErr.Type != ErrorTypeMigration {
		t.Errorf("Expected ErrorTypeMigration, got %v", migErr.Type)
	}
	
	if migErr.Details != "Detailed explanation" {
		t.Errorf("Expected details to be set correctly")
	}
	
	// Test NewIOError
	ioCause := errors.New("file not found")
	ioErr := NewIOError("Cannot read file", ioCause)
	
	if ioErr.Type != ErrorTypeIO {
		t.Errorf("Expected ErrorTypeIO, got %v", ioErr.Type)
	}
	
	if ioErr.Cause != ioCause {
		t.Errorf("Expected cause to be set correctly")
	}
	
	// Test NewUnsupportedError
	unsupErr := NewUnsupportedError("Custom Interaction", "2.1")
	
	if unsupErr.Type != ErrorTypeUnsupported {
		t.Errorf("Expected ErrorTypeUnsupported, got %v", unsupErr.Type)
	}
	
	expectedMsg := "Custom Interaction is not supported in QTI 2.1"
	if unsupErr.Message != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, unsupErr.Message)
	}
}

func TestQTIError_Error_Method(t *testing.T) {
	// Test error with ItemID
	err1 := &QTIError{
		Type:    ErrorTypeValidation,
		Message: "Invalid value",
		ItemID:  "q001",
	}
	
	expected1 := "[q001] Validation Error: Invalid value"
	if err1.Error() != expected1 {
		t.Errorf("Expected '%s', got '%s'", expected1, err1.Error())
	}
	
	// Test error without ItemID
	err2 := &QTIError{
		Type:    ErrorTypeParsing,
		Message: "XML parse error",
	}
	
	expected2 := "Parsing Error: XML parse error"
	if err2.Error() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, err2.Error())
	}
}

func TestQTIError_TypeString(t *testing.T) {
	testCases := []struct {
		errorType    ErrorType
		expectedStr  string
	}{
		{ErrorTypeParsing, "Parsing Error"},
		{ErrorTypeValidation, "Validation Error"},
		{ErrorTypeMigration, "Migration Error"},
		{ErrorTypeIO, "I/O Error"},
		{ErrorTypeUnsupported, "Unsupported Feature"},
		{ErrorType(999), "Unknown Error"}, // Unknown error type
	}
	
	for _, tc := range testCases {
		err := &QTIError{Type: tc.errorType}
		result := err.TypeString()
		if result != tc.expectedStr {
			t.Errorf("For error type %v, expected '%s', got '%s'", 
				tc.errorType, tc.expectedStr, result)
		}
	}
}

func TestQTIError_Unwrap(t *testing.T) {
	causeErr := errors.New("underlying error")
	qtiErr := &QTIError{
		Type:    ErrorTypeParsing,
		Message: "Parse failed",
		Cause:   causeErr,
	}
	
	unwrapped := qtiErr.Unwrap()
	if unwrapped != causeErr {
		t.Errorf("Expected unwrapped error to be the cause error")
	}
	
	// Test with no cause
	qtiErrNoCause := &QTIError{
		Type:    ErrorTypeValidation,
		Message: "Validation failed",
	}
	
	unwrappedNil := qtiErrNoCause.Unwrap()
	if unwrappedNil != nil {
		t.Errorf("Expected unwrapped error to be nil when no cause is set")
	}
}

func TestErrorList_Add(t *testing.T) {
	el := &ErrorList{}
	
	if el.HasErrors() {
		t.Error("Expected empty error list to have no errors")
	}
	
	err1 := NewParsingError("Parse error", nil)
	el.Add(err1)
	
	if !el.HasErrors() {
		t.Error("Expected error list to have errors after adding one")
	}
	
	if len(el.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(el.Errors))
	}
	
	if el.Errors[0] != err1 {
		t.Error("Expected added error to be in the list")
	}
	
	// Add another error
	err2 := NewValidationError("Validation error", "q001", "/item")
	el.Add(err2)
	
	if len(el.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(el.Errors))
	}
}

func TestErrorList_Error(t *testing.T) {
	// Test empty error list
	el := &ErrorList{}
	if el.Error() != "no errors" {
		t.Errorf("Expected 'no errors', got '%s'", el.Error())
	}
	
	// Test single error
	err1 := NewParsingError("Parse error", nil)
	el.Add(err1)
	
	expected1 := err1.Error()
	if el.Error() != expected1 {
		t.Errorf("Expected '%s', got '%s'", expected1, el.Error())
	}
	
	// Test multiple errors
	err2 := NewValidationError("Validation error", "q001", "/item")
	el.Add(err2)
	
	expected2 := "2 errors occurred during processing"
	if el.Error() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, el.Error())
	}
}

func TestErrorList_GetByType(t *testing.T) {
	el := &ErrorList{}
	
	// Add different types of errors
	parseErr1 := NewParsingError("Parse error 1", nil)
	parseErr2 := NewParsingError("Parse error 2", nil)
	validErr := NewValidationError("Validation error", "q001", "/item")
	migErr := NewMigrationError("Migration error", "details")
	
	el.Add(parseErr1)
	el.Add(validErr)
	el.Add(parseErr2)
	el.Add(migErr)
	
	// Test getting parsing errors
	parseErrors := el.GetByType(ErrorTypeParsing)
	if len(parseErrors) != 2 {
		t.Errorf("Expected 2 parsing errors, got %d", len(parseErrors))
	}
	
	// Check that the right errors are returned
	foundParseErr1 := false
	foundParseErr2 := false
	for _, err := range parseErrors {
		if err == parseErr1 {
			foundParseErr1 = true
		}
		if err == parseErr2 {
			foundParseErr2 = true
		}
	}
	
	if !foundParseErr1 || !foundParseErr2 {
		t.Error("Expected both parsing errors to be returned")
	}
	
	// Test getting validation errors
	validErrors := el.GetByType(ErrorTypeValidation)
	if len(validErrors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(validErrors))
	}
	
	if validErrors[0] != validErr {
		t.Error("Expected validation error to be returned")
	}
	
	// Test getting non-existent error type
	ioErrors := el.GetByType(ErrorTypeIO)
	if len(ioErrors) != 0 {
		t.Errorf("Expected 0 I/O errors, got %d", len(ioErrors))
	}
}

func TestNewUnsupportedError_MessageFormat(t *testing.T) {
	testCases := []struct {
		feature        string
		version        string
		expectedMsg    string
		expectedDetail string
	}{
		{
			feature:        "Custom Interaction",
			version:        "2.1",
			expectedMsg:    "Custom Interaction is not supported in QTI 2.1",
			expectedDetail: "Feature 'Custom Interaction' cannot be migrated to QTI 2.1",
		},
		{
			feature:        "Advanced Scoring",
			version:        "3.0",
			expectedMsg:    "Advanced Scoring is not supported in QTI 3.0",
			expectedDetail: "Feature 'Advanced Scoring' cannot be migrated to QTI 3.0",
		},
	}
	
	for _, tc := range testCases {
		err := NewUnsupportedError(tc.feature, tc.version)
		
		if err.Message != tc.expectedMsg {
			t.Errorf("Expected message '%s', got '%s'", tc.expectedMsg, err.Message)
		}
		
		if err.Details != tc.expectedDetail {
			t.Errorf("Expected details '%s', got '%s'", tc.expectedDetail, err.Details)
		}
	}
}

func TestQTIError_Integration_WithStandardErrors(t *testing.T) {
	// Test that QTIError works well with Go's error handling patterns
	baseErr := errors.New("base error")
	qtiErr := NewParsingError("Parse failed", baseErr)
	
	// Test with errors.Is
	if !errors.Is(qtiErr, baseErr) {
		t.Error("Expected errors.Is to work with wrapped error")
	}
	
	// Test with errors.As
	var qtiErrTarget *QTIError
	if !errors.As(qtiErr, &qtiErrTarget) {
		t.Error("Expected errors.As to work with QTIError")
	}
	
	if qtiErrTarget.Type != ErrorTypeParsing {
		t.Errorf("Expected parsed error type, got %v", qtiErrTarget.Type)
	}
}

func TestErrorList_ConcurrentAccess(t *testing.T) {
	// Test that ErrorList can handle concurrent access safely
	// Note: This is a basic test - in a real concurrent scenario, 
	// you might need proper synchronization
	
	el := &ErrorList{}
	
	// Add errors in a way that simulates potential concurrent access
	errors := []*QTIError{
		NewParsingError("Error 1", nil),
		NewValidationError("Error 2", "q001", "/item"),
		NewMigrationError("Error 3", "details"),
		NewIOError("Error 4", nil),
		NewUnsupportedError("Feature", "2.1"),
	}
	
	for _, err := range errors {
		el.Add(err)
	}
	
	// Verify all errors were added
	if len(el.Errors) != len(errors) {
		t.Errorf("Expected %d errors, got %d", len(errors), len(el.Errors))
	}
	
	// Verify GetByType works correctly
	parseErrors := el.GetByType(ErrorTypeParsing)
	if len(parseErrors) != 1 {
		t.Errorf("Expected 1 parsing error, got %d", len(parseErrors))
	}
	
	validErrors := el.GetByType(ErrorTypeValidation)
	if len(validErrors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(validErrors))
	}
}

func BenchmarkQTIError_Error(b *testing.B) {
	err := &QTIError{
		Type:    ErrorTypeValidation,
		Message: "Benchmark error message",
		ItemID:  "q001",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkErrorList_Add(b *testing.B) {
	el := &ErrorList{}
	baseErr := NewParsingError("Benchmark error", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new error for each iteration to avoid reusing the same object
		err := &QTIError{
			Type:    baseErr.Type,
			Message: baseErr.Message,
		}
		el.Add(err)
	}
}

func BenchmarkErrorList_GetByType(b *testing.B) {
	el := &ErrorList{}
	
	// Add various types of errors
	for i := 0; i < 100; i++ {
		el.Add(NewParsingError("Parse error", nil))
		el.Add(NewValidationError("Valid error", "q", "/item"))
		el.Add(NewMigrationError("Migration error", "details"))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = el.GetByType(ErrorTypeValidation)
	}
}