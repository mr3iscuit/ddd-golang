package event

import (
	"time"

	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoCompletedEvent represents a domain event when a Todo is completed
type TodoCompletedEvent struct {
	TodoID      model.TodoID
	CompletedAt time.Time
}

// NewTodoCompletedEvent creates a new TodoCompletedEvent
func NewTodoCompletedEvent(todoID model.TodoID) *TodoCompletedEvent {
	return &TodoCompletedEvent{
		TodoID:      todoID,
		CompletedAt: time.Now(),
	}
}
