package service

import (
	"strings"

	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoDomainService handles domain-specific business logic for todos
// Implements port.TodoDomainServicePort
type TodoDomainService struct{}

// Ensure TodoDomainService implements TodoDomainServicePort
var _ port.TodoDomainServicePort = (*TodoDomainService)(nil)

// NewTodoDomainService creates a new todo domain service
func NewTodoDomainService() *TodoDomainService {
	return &TodoDomainService{}
}

// ValidateTitle validates a todo title
func (s *TodoDomainService) ValidateTitle(title string) *model.DomainError {
	if strings.TrimSpace(title) == "" {
		return model.ErrEmptyTitle
	}
	if len(title) > 100 {
		return model.ErrTitleTooLong
	}
	return nil
}

// ValidateDescription validates a todo description
func (s *TodoDomainService) ValidateDescription(description string) *model.DomainError {
	if len(description) > 1000 {
		return model.ErrInvalidDescription
	}
	return nil
}

// ValidatePriority validates a todo priority
func (s *TodoDomainService) ValidatePriority(priority string) *model.DomainError {
	switch priority {
	case "low", "medium", "high":
		return nil
	default:
		return model.ErrInvalidPriority
	}
}

// ValidateCreateTodoCommand validates all fields for creating a todo
func (s *TodoDomainService) ValidateCreateTodoCommand(title string, description string, priority string) *model.DomainError {
	if err := s.ValidateTitle(title); err != nil {
		return err
	}
	if err := s.ValidateDescription(description); err != nil {
		return err
	}
	if err := s.ValidatePriority(priority); err != nil {
		return err
	}
	return nil
}

// ValidateUpdateTodoCommand validates all fields for updating a todo
func (s *TodoDomainService) ValidateUpdateTodoCommand(title string, description string, priority string) *model.DomainError {
	if title != "" {
		if err := s.ValidateTitle(title); err != nil {
			return err
		}
	}
	if description != "" {
		if err := s.ValidateDescription(description); err != nil {
			return err
		}
	}
	if priority != "" {
		if err := s.ValidatePriority(priority); err != nil {
			return err
		}
	}
	return nil
}
