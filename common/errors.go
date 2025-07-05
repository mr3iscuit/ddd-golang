package shared

import (
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// CommonError creates a common domain error
func CommonError(message string, reason string) *model.DomainError {
	return model.NewDomainError(500, 500, message, reason, nil)
}

// ValidationError creates a validation domain error
func ValidationError(message string, details map[string]string) *model.DomainError {
	return model.NewDomainError(400, 400, message, "Validation failed", details)
}

// NotFoundError creates a not found domain error
func NotFoundError(resource string, id string) *model.DomainError {
	return model.NewDomainError(404, 404,
		resource+" not found",
		resource+" with id "+id+" not found",
		map[string]string{"resource": resource, "id": id})
}
