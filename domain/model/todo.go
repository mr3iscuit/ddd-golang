package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TodoID represents a unique Todo identifier following DDD naming
type TodoID string

// TodoStatus represents the completion status of a Todo
type TodoStatus string

const (
	TodoStatusPending   TodoStatus = "pending"
	TodoStatusCompleted TodoStatus = "completed"
	TodoStatusArchived  TodoStatus = "archived"
)

// TodoPriority represents the priority level of a Todo
type TodoPriority string

const (
	TodoPriorityLow    TodoPriority = "low"
	TodoPriorityMedium TodoPriority = "medium"
	TodoPriorityHigh   TodoPriority = "high"
)

// Todo represents the Todo aggregate root in DDD
type Todo struct {
	id          TodoID
	title       string
	description string
	status      TodoStatus
	priority    TodoPriority
	createdAt   time.Time
	updatedAt   time.Time
	completedAt *time.Time
}

// NewTodo creates a new Todo aggregate root with descriptive factory method
func NewTodo(title string, description string, priority TodoPriority) *Todo {
	now := time.Now()
	return &Todo{
		id:          TodoID(uuid.NewString()),
		title:       title,
		description: description,
		status:      TodoStatusPending,
		priority:    priority,
		createdAt:   now,
		updatedAt:   now,
		completedAt: nil,
	}
}

func NewTodoWithAllFields(
	id TodoID,
	title string,
	description string,
	status TodoStatus,
	priority TodoPriority,
	createdAt time.Time,
	updatedAt time.Time,
	completedAt *time.Time,
) *Todo {
	return &Todo{
		id:          id,
		title:       title,
		description: description,
		status:      TodoStatusPending,
		priority:    priority,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		completedAt: completedAt,
	}
}

// NewSimpleTodo creates a new Todo with minimal required fields
func NewSimpleTodo(title string) *Todo {
	return NewTodo(title, "", TodoPriorityMedium)
}

// NewTodoFromData reconstructs a Todo object from persistent data
func NewTodoFromData(id TodoID, title, description string, status TodoStatus, priority TodoPriority, createdAt, updatedAt time.Time, completedAt *time.Time) *Todo {
	return &Todo{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		priority:    priority,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		completedAt: completedAt,
	}
}

// Getters following DDD encapsulation principles with descriptive names
func (t *Todo) GetID() TodoID {
	return t.id
}

func (t *Todo) GetTitle() string {
	return t.title
}

func (t *Todo) GetDescription() string {
	return t.description
}

func (t *Todo) GetStatus() TodoStatus {
	return t.status
}

func (t *Todo) GetPriority() TodoPriority {
	return t.priority
}

func (t *Todo) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Todo) GetUpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Todo) GetCompletedAt() *time.Time {
	return t.completedAt
}

// IsCompleted checks if the todo is completed
func (t *Todo) IsCompleted() bool {
	return t.status == TodoStatusCompleted
}

// IsPending checks if the todo is pending
func (t *Todo) IsPending() bool {
	return t.status == TodoStatusPending
}

// IsArchived checks if the todo is archived
func (t *Todo) IsArchived() bool {
	return t.status == TodoStatusArchived
}

// MarkAsCompleted is a domain behavior that enforces business rules
func (t *Todo) MarkAsCompleted() error {
	if t.IsCompleted() {
		return errors.New("todo is already completed")
	}
	if t.IsArchived() {
		return errors.New("cannot complete an archived todo")
	}

	now := time.Now()
	t.status = TodoStatusCompleted
	t.completedAt = &now
	t.updatedAt = now
	return nil
}

// MarkAsPending resets the todo to pending status
func (t *Todo) MarkAsPending() error {
	if t.IsCompleted() {
		return errors.New("cannot mark completed todo as pending")
	}

	t.status = TodoStatusPending
	t.completedAt = nil
	t.updatedAt = time.Now()
	return nil
}

// ArchiveTodo archives the todo
func (t *Todo) ArchiveTodo() error {
	if t.IsArchived() {
		return errors.New("todo is already archived")
	}

	t.status = TodoStatusArchived
	t.updatedAt = time.Now()
	return nil
}

// UpdateTitle allows updating the todo title with validation
func (t *Todo) UpdateTitle(newTitle string) error {
	if newTitle == "" {
		return errors.New("title cannot be empty")
	}
	if len(newTitle) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	t.title = newTitle
	t.updatedAt = time.Now()
	return nil
}

// UpdateDescription allows updating the todo description
func (t *Todo) UpdateDescription(newDescription string) error {
	if len(newDescription) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	t.description = newDescription
	t.updatedAt = time.Now()
	return nil
}

// UpdatePriority allows updating the todo priority
func (t *Todo) UpdatePriority(newPriority TodoPriority) error {
	switch newPriority {
	case TodoPriorityLow, TodoPriorityMedium, TodoPriorityHigh:
		t.priority = newPriority
		t.updatedAt = time.Now()
		return nil
	default:
		return errors.New("invalid priority level")
	}
}

// GetElapsedTimeSinceCreation returns the time elapsed since todo creation
func (t *Todo) GetElapsedTimeSinceCreation() time.Duration {
	return time.Since(t.createdAt)
}

// GetElapsedTimeSinceCompletion returns the time elapsed since completion (if completed)
func (t *Todo) GetElapsedTimeSinceCompletion() (time.Duration, error) {
	if !t.IsCompleted() || t.completedAt == nil {
		return 0, errors.New("todo is not completed")
	}
	return time.Since(*t.completedAt), nil
}
