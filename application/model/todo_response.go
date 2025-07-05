package model

import (
	"time"

	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoResponse represents a todo item in the application layer
type TodoResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"created-at"`
	CompletedAt *time.Time `json:"completed-at,omitempty"`
}

// TodoListResponse represents a list of todos
type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
	Count int            `json:"count"`
}

// TodoResponseMapper maps a domain Todo to a TodoResponse
func TodoResponseMapper(todo *model.Todo) TodoResponse {
	response := TodoResponse{
		ID:          string(todo.GetID()),
		Title:       todo.GetTitle(),
		Description: todo.GetDescription(),
		Status:      string(todo.GetStatus()),
		Priority:    string(todo.GetPriority()),
		CreatedAt:   todo.GetCreatedAt(),
	}

	if todo.GetCompletedAt() != nil {
		response.CompletedAt = todo.GetCompletedAt()
	}

	return response
}

// TodoListResponseMapper maps a slice of domain Todos to a TodoListResponse
func TodoListResponseMapper(todos []*model.Todo) TodoListResponse {
	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = TodoResponseMapper(todo)
	}

	return TodoListResponse{
		Todos: responses,
		Count: len(responses),
	}
}
