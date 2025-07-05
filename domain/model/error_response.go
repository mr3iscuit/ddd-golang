package model

// DomainErrorResponse represents the error response structure as part of the domain contract
type DomainErrorResponse struct {
	ErrorCode      int               `json:"error-code"`
	HttpStatus     int               `json:"http-status"`
	ErrorMessage   string            `json:"error-message"`
	InternalReason string            `json:"internal-reason,omitempty"`
	Details        map[string]string `json:"details,omitempty"`
}

// ToResponse maps a domain error to its response representation
func (e *DomainError) ToResponse() DomainErrorResponse {
	return DomainErrorResponse{
		ErrorCode:      e.GetErrorCode(),
		HttpStatus:     e.GetHttpStatus(),
		ErrorMessage:   e.GetErrorMessage(),
		InternalReason: e.GetInternalReason(),
		Details:        e.GetDetails(),
	}
}

// ToResponseWithInternal maps a domain error to its response representation with optional internal details
func (e *DomainError) ToResponseWithInternal(includeInternal bool) DomainErrorResponse {
	response := DomainErrorResponse{
		ErrorCode:    e.GetErrorCode(),
		HttpStatus:   e.GetHttpStatus(),
		ErrorMessage: e.GetErrorMessage(),
		Details:      e.GetDetails(),
	}

	if includeInternal {
		response.InternalReason = e.GetInternalReason()
	}

	return response
}
