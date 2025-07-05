package model

import (
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// ErrorResponse is now an alias to the domain's error response structure
type ErrorResponse = model.DomainErrorResponse

// ErrorResponseMapper maps a domain error to an error response using domain's mapping
func ErrorResponseMapper(domainError model.DomainErrorPort) ErrorResponse {
	if de, ok := domainError.(*model.DomainError); ok {
		return de.ToResponse()
	}
	// Fallback for other error types
	return ErrorResponse{
		ErrorCode:      9999,
		HttpStatus:     500,
		ErrorMessage:   domainError.Error(),
		InternalReason: "Unknown error type",
	}
}

// ErrorResponseMapperWithInternal maps a domain error to an error response including internal details
func ErrorResponseMapperWithInternal(domainError model.DomainErrorPort, includeInternal bool) ErrorResponse {
	if de, ok := domainError.(*model.DomainError); ok {
		return de.ToResponseWithInternal(includeInternal)
	}
	// Fallback for other error types
	response := ErrorResponse{
		ErrorCode:    9999,
		HttpStatus:   500,
		ErrorMessage: domainError.Error(),
	}
	if includeInternal {
		response.InternalReason = "Unknown error type"
	}
	return response
}
