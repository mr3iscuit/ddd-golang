package port

import "github.com/mr3iscuit/ddd-golang/domain/model"

// TodoDomainServicePort defines the interface for domain service operations
type TodoDomainServicePort interface {
	ValidateTitle(title string) *model.DomainError
	ValidateDescription(description string) *model.DomainError
	ValidatePriority(priority string) *model.DomainError
	ValidateCreateTodoCommand(title string, description string, priority string) *model.DomainError
	ValidateUpdateTodoCommand(title string, description string, priority string) *model.DomainError
}
