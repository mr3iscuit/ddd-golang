package service

import (
	"strings"

	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoDomainService handles domain-specific business logic for todos
type TodoDomainService struct{}

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
