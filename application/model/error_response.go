package model

import (
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// ErrorResponse represents the error response structure for HTTP responses
type ErrorResponse struct {
	ErrorCode      int               `json:"error-code"`
	HttpStatus     int               `json:"http-status"`
	ErrorMessage   string            `json:"error-message"`
	InternalReason string            `json:"internal-reason,omitempty"`
	Details        map[string]string `json:"details,omitempty"`
}

// ErrorResponseMapper maps a domain error to an error response DTO
func ErrorResponseMapper(domainError model.DomainErrorPort) ErrorResponse {
	return ErrorResponse{
		ErrorCode:      domainError.GetErrorCode(),
		HttpStatus:     domainError.GetHttpStatus(),
		ErrorMessage:   domainError.GetErrorMessage(),
		InternalReason: domainError.GetInternalReason(),
		Details:        domainError.GetDetails(),
	}
}

// ErrorResponseMapperWithInternal maps a domain error to an error response DTO including internal details
func ErrorResponseMapperWithInternal(domainError model.DomainErrorPort, includeInternal bool) ErrorResponse {
	response := ErrorResponse{
		ErrorCode:    domainError.GetErrorCode(),
		HttpStatus:   domainError.GetHttpStatus(),
		ErrorMessage: domainError.GetErrorMessage(),
		Details:      domainError.GetDetails(),
	}

	if includeInternal {
		response.InternalReason = domainError.GetInternalReason()
	}

	return response
}
