package port

import "github.com/mr3iscuit/ddd-golang/domain/model"

// TodoRepositoryPort is the outbound port for Todo persistence
// (previously domain/repository.TodoRepository)
type TodoRepositoryPort interface {
	Save(todo *model.Todo) error
	FindByID(id model.TodoID) (*model.Todo, error)
	FindAll() ([]*model.Todo, error)
	Delete(id model.TodoID) error
}
