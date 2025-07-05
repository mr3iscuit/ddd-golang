package model

// DomainError represents a domain-specific error following DDD principles
type DomainError struct {
	errorCode      int
	httpStatus     int
	errorMessage   string
	internalReason string
	details        map[string]string
}

// DomainErrorPort defines the interface for domain errors
type DomainErrorPort interface {
	GetErrorCode() int
	GetHttpStatus() int
	GetErrorMessage() string
	GetInternalReason() string
	GetDetails() map[string]string
	Error() string
}

// GetErrorCode returns the error code
func (e *DomainError) GetErrorCode() int {
	return e.errorCode
}

// GetHttpStatus returns the HTTP status
func (e *DomainError) GetHttpStatus() int {
	return e.httpStatus
}

// GetErrorMessage returns the error message
func (e *DomainError) GetErrorMessage() string {
	return e.errorMessage
}

// GetInternalReason returns the internal reason
func (e *DomainError) GetInternalReason() string {
	return e.internalReason
}

// GetDetails returns the error details
func (e *DomainError) GetDetails() map[string]string {
	return e.details
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.errorMessage
}

// Predefined domain errors organized by category

// Validation errors (1000-1999)
var (
	ErrInvalidTitle = &DomainError{
		errorCode:      1001,
		httpStatus:     400,
		errorMessage:   "Invalid title",
		internalReason: "Title validation failed",
		details:        nil,
	}

	ErrInvalidDescription = &DomainError{
		errorCode:      1002,
		httpStatus:     400,
		errorMessage:   "Invalid description",
		internalReason: "Description validation failed",
		details:        nil,
	}

	ErrInvalidPriority = &DomainError{
		errorCode:      1003,
		httpStatus:     400,
		errorMessage:   "Invalid priority",
		internalReason: "Priority must be low, medium, or high",
		details:        nil,
	}

	ErrEmptyTitle = &DomainError{
		errorCode:      1004,
		httpStatus:     400,
		errorMessage:   "Title cannot be empty",
		internalReason: "Empty title provided",
		details:        nil,
	}

	ErrTitleTooLong = &DomainError{
		errorCode:      1005,
		httpStatus:     400,
		errorMessage:   "Title too long",
		internalReason: "Title exceeds maximum length of 100 characters",
		details:        map[string]string{"max_length": "100"},
	}
)

// Not found errors (2000-2999)
var (
	ErrTodoNotFound = &DomainError{
		errorCode:      2001,
		httpStatus:     404,
		errorMessage:   "Todo not found",
		internalReason: "Todo with specified ID not found",
		details:        nil,
	}
)

// Operation errors (3000-3999)
var (
	ErrCannotCompleteTodo = &DomainError{
		errorCode:      3001,
		httpStatus:     400,
		errorMessage:   "Cannot complete todo",
		internalReason: "Todo cannot be completed",
		details:        nil,
	}

	ErrCannotArchiveTodo = &DomainError{
		errorCode:      3002,
		httpStatus:     400,
		errorMessage:   "Cannot archive todo",
		internalReason: "Todo cannot be archived",
		details:        nil,
	}
)

// Repository errors (4000-4999)
var (
	ErrRepositoryNotInitialized = &DomainError{
		errorCode:      4001,
		httpStatus:     500,
		errorMessage:   "Repository not initialized",
		internalReason: "Repository is nil",
		details:        map[string]string{"operation": "list_todos"},
	}

	ErrFailedToSaveTodo = &DomainError{
		errorCode:      4002,
		httpStatus:     500,
		errorMessage:   "Failed to save todo",
		internalReason: "Database save operation failed",
		details:        nil,
	}

	ErrFailedToSaveCompletedTodo = &DomainError{
		errorCode:      4003,
		httpStatus:     500,
		errorMessage:   "Failed to save completed todo",
		internalReason: "Database save operation failed for completed todo",
		details:        nil,
	}

	ErrFailedToSaveArchivedTodo = &DomainError{
		errorCode:      4004,
		httpStatus:     500,
		errorMessage:   "Failed to save archived todo",
		internalReason: "Database save operation failed for archived todo",
		details:        nil,
	}

	ErrFailedToRetrieveTodos = &DomainError{
		errorCode:      4005,
		httpStatus:     500,
		errorMessage:   "Failed to retrieve todos",
		internalReason: "Database retrieve operation failed",
		details:        map[string]string{"operation": "list_todos"},
	}
)

// HTTP errors (5000-5999)
var (
	ErrInvalidJSON = &DomainError{
		errorCode:      5001,
		httpStatus:     400,
		errorMessage:   "Invalid JSON",
		internalReason: "JSON parsing failed",
		details:        nil,
	}
)

// Test errors (9000-9999)
var (
	ErrTestError = &DomainError{
		errorCode:      9001,
		httpStatus:     400,
		errorMessage:   "Test error message",
		internalReason: "This is a test error for testing error handling",
		details:        map[string]string{"test": "true"},
	}
)

// NewDomainError creates a new domain error (kept for backward compatibility)
func NewDomainError(
	errorCode int,
	httpStatus int,
	errorMessage string,
	internalReason string,
	details map[string]string,
) *DomainError {
	return &DomainError{
		errorCode:      errorCode,
		httpStatus:     httpStatus,
		errorMessage:   errorMessage,
		internalReason: internalReason,
		details:        details,
	}
}
