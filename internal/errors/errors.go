package errors

import (
	"fmt"
)

type ErrorType int

const (
	ErrorTypeParsing ErrorType = iota
	ErrorTypeValidation
	ErrorTypeMigration
	ErrorTypeIO
	ErrorTypeUnsupported
)

type QTIError struct {
	Type        ErrorType
	Message     string
	Details     string
	ItemID      string
	ElementPath string
	Cause       error
}

func (e *QTIError) Error() string {
	if e.ItemID != "" {
		return fmt.Sprintf("[%s] %s: %s", e.ItemID, e.TypeString(), e.Message)
	}
	return fmt.Sprintf("%s: %s", e.TypeString(), e.Message)
}

func (e *QTIError) TypeString() string {
	switch e.Type {
	case ErrorTypeParsing:
		return "Parsing Error"
	case ErrorTypeValidation:
		return "Validation Error"
	case ErrorTypeMigration:
		return "Migration Error"
	case ErrorTypeIO:
		return "I/O Error"
	case ErrorTypeUnsupported:
		return "Unsupported Feature"
	default:
		return "Unknown Error"
	}
}

func (e *QTIError) Unwrap() error {
	return e.Cause
}

func NewParsingError(message string, cause error) *QTIError {
	return &QTIError{
		Type:    ErrorTypeParsing,
		Message: message,
		Cause:   cause,
	}
}

func NewValidationError(message, itemID, elementPath string) *QTIError {
	return &QTIError{
		Type:        ErrorTypeValidation,
		Message:     message,
		ItemID:      itemID,
		ElementPath: elementPath,
	}
}

func NewMigrationError(message, details string) *QTIError {
	return &QTIError{
		Type:    ErrorTypeMigration,
		Message: message,
		Details: details,
	}
}

func NewIOError(message string, cause error) *QTIError {
	return &QTIError{
		Type:    ErrorTypeIO,
		Message: message,
		Cause:   cause,
	}
}

func NewUnsupportedError(feature, version string) *QTIError {
	return &QTIError{
		Type:    ErrorTypeUnsupported,
		Message: fmt.Sprintf("%s is not supported in QTI %s", feature, version),
		Details: fmt.Sprintf("Feature '%s' cannot be migrated to QTI %s", feature, version),
	}
}

type ErrorList struct {
	Errors []*QTIError
}

func (el *ErrorList) Add(err *QTIError) {
	el.Errors = append(el.Errors, err)
}

func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

func (el *ErrorList) Error() string {
	if len(el.Errors) == 0 {
		return "no errors"
	}
	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred during processing", len(el.Errors))
}

func (el *ErrorList) GetByType(errorType ErrorType) []*QTIError {
	var filtered []*QTIError
	for _, err := range el.Errors {
		if err.Type == errorType {
			filtered = append(filtered, err)
		}
	}
	return filtered
}