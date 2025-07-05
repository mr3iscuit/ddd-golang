package model

// DomainError represents a domain-specific error following DDD principles
type DomainError struct {
	statusCode     int
	httpStatus     int
	errorMessage   string
	internalReason string
	details        map[string]string
}

// DomainErrorPort defines the interface for domain errors
type DomainErrorPort interface {
	GetStatusCode() int
	GetHttpStatus() int
	GetErrorMessage() string
	GetInternalReason() string
	GetDetails() map[string]string
	Error() string
}

// GetStatusCode returns the status code
func (e *DomainError) GetStatusCode() int {
	return e.statusCode
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

// Predefined domain errors
var (
	// Validation errors
	ErrInvalidTitle = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Invalid title",
		internalReason: "Title validation failed",
		details:        nil,
	}

	ErrInvalidDescription = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Invalid description",
		internalReason: "Description validation failed",
		details:        nil,
	}

	ErrInvalidPriority = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Invalid priority",
		internalReason: "Priority must be low, medium, or high",
		details:        nil,
	}

	ErrEmptyTitle = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Title cannot be empty",
		internalReason: "Empty title provided",
		details:        nil,
	}

	ErrTitleTooLong = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Title too long",
		internalReason: "Title exceeds maximum length of 100 characters",
		details:        map[string]string{"max_length": "100"},
	}

	// Not found errors
	ErrTodoNotFound = &DomainError{
		statusCode:     404,
		httpStatus:     404,
		errorMessage:   "Todo not found",
		internalReason: "Todo with specified ID not found",
		details:        nil,
	}

	// Operation errors
	ErrCannotCompleteTodo = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Cannot complete todo",
		internalReason: "Todo cannot be completed",
		details:        nil,
	}

	ErrCannotArchiveTodo = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Cannot archive todo",
		internalReason: "Todo cannot be archived",
		details:        nil,
	}

	// Repository errors
	ErrRepositoryNotInitialized = &DomainError{
		statusCode:     500,
		httpStatus:     500,
		errorMessage:   "Repository not initialized",
		internalReason: "Repository is nil",
		details:        map[string]string{"operation": "list_todos"},
	}

	ErrFailedToSaveTodo = &DomainError{
		statusCode:     500,
		httpStatus:     500,
		errorMessage:   "Failed to save todo",
		internalReason: "Database save operation failed",
		details:        nil,
	}

	ErrFailedToSaveCompletedTodo = &DomainError{
		statusCode:     500,
		httpStatus:     500,
		errorMessage:   "Failed to save completed todo",
		internalReason: "Database save operation failed for completed todo",
		details:        nil,
	}

	ErrFailedToSaveArchivedTodo = &DomainError{
		statusCode:     500,
		httpStatus:     500,
		errorMessage:   "Failed to save archived todo",
		internalReason: "Database save operation failed for archived todo",
		details:        nil,
	}

	ErrFailedToRetrieveTodos = &DomainError{
		statusCode:     500,
		httpStatus:     500,
		errorMessage:   "Failed to retrieve todos",
		internalReason: "Database retrieve operation failed",
		details:        map[string]string{"operation": "list_todos"},
	}

	// HTTP errors
	ErrInvalidJSON = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Invalid JSON",
		internalReason: "JSON parsing failed",
		details:        nil,
	}

	// Test error
	ErrTestError = &DomainError{
		statusCode:     400,
		httpStatus:     400,
		errorMessage:   "Test error message",
		internalReason: "This is a test error for testing error handling",
		details:        map[string]string{"test": "true"},
	}
)

// NewDomainError creates a new domain error (kept for backward compatibility)
func NewDomainError(
	statusCode int,
	httpStatus int,
	errorMessage string,
	internalReason string,
	details map[string]string,
) *DomainError {
	return &DomainError{
		statusCode:     statusCode,
		httpStatus:     httpStatus,
		errorMessage:   errorMessage,
		internalReason: internalReason,
		details:        details,
	}
}
