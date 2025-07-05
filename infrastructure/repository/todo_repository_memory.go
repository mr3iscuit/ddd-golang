package repository

import (
	"fmt"
	"sync"

	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// InMemoryTodoRepository implements port.TodoRepositoryPort
// (was domain/repository.TodoRepository)
type InMemoryTodoRepository struct {
	todos map[model.TodoID]*model.Todo
	mutex sync.RWMutex
}

// NewInMemoryTodoRepository creates a new in-memory Todo repository
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[model.TodoID]*model.Todo),
	}
}

var _ port.TodoRepositoryPort = (*InMemoryTodoRepository)(nil)

// Save stores a Todo in memory
func (r *InMemoryTodoRepository) Save(todo *model.Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.todos[todo.GetID()] = todo
	return nil
}

// FindByID retrieves a Todo by ID
func (r *InMemoryTodoRepository) FindByID(id model.TodoID) (*model.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	todo, exists := r.todos[id]
	if !exists {
		return nil, fmt.Errorf("todo with id %s not found", id)
	}
	return todo, nil
}

// FindAll retrieves all Todos
func (r *InMemoryTodoRepository) FindAll() ([]*model.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var todos []*model.Todo
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

// Delete removes a Todo by ID
func (r *InMemoryTodoRepository) Delete(id model.TodoID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.todos[id]; !exists {
		return fmt.Errorf("todo with id %s not found", id)
	}
	delete(r.todos, id)
	return nil
}
