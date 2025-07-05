package postgres

import "github.com/mr3iscuit/ddd-golang/domain/model"

func fromModel(todo *model.Todo) *TodoRecord {
	return &TodoRecord{
		ID:          string(todo.GetID()),
		Title:       todo.GetTitle(),
		Description: todo.GetDescription(),
		Priority:    string(todo.GetPriority()),
		Status:      string(todo.GetStatus()),
		CreatedAt:   todo.GetCreatedAt(),
		UpdatedAt:   todo.GetUpdatedAt(),
		CompletedAt: todo.GetCompletedAt(),
	}
}

func toModel(r *TodoRecord) *model.Todo {
	return model.NewTodoFromData(
		model.TodoID(r.ID),
		r.Title,
		r.Description,
		model.TodoStatus(r.Status),
		model.TodoPriority(r.Priority),
		r.CreatedAt,
		r.UpdatedAt,
		r.CompletedAt,
	)
}
