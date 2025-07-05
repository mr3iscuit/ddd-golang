package port

import (
	"github.com/mr3iscuit/ddd-golang/application/command"
	appmodel "github.com/mr3iscuit/ddd-golang/application/model"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoUseCasePort defines the inbound port for Todo use cases
type TodoUseCasePort interface {
	CreateTodoUseCase(cmd command.CreateTodoCommand) (model.TodoID, *model.DomainError)
	UpdateTodoUseCase(cmd command.UpdateTodoCommand) *model.DomainError
	CompleteTodoUseCase(id model.TodoID) *model.DomainError
	ArchiveTodoUseCase(id model.TodoID) *model.DomainError
	GetTodoUseCase(id model.TodoID) (*appmodel.TodoResponse, *model.DomainError)
	ListTodosUseCase() (*appmodel.TodoListResponse, *model.DomainError)
	TestErrorUseCase() *model.DomainError
}
