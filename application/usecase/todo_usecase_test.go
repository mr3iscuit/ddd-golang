package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mr3iscuit/ddd-golang/application/command"
	"github.com/mr3iscuit/ddd-golang/domain/model"
	"github.com/mr3iscuit/ddd-golang/domain/service"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) Save(todo *model.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) FindByID(id model.TodoID) (*model.Todo, error) {
	args := m.Called(id)
	if todo, ok := args.Get(0).(*model.Todo); ok {
		return todo, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTodoRepository) FindAll() ([]*model.Todo, error) {
	args := m.Called()
	if todos, ok := args.Get(0).([]*model.Todo); ok {
		return todos, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTodoRepository) Delete(id model.TodoID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTodoUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	cmd := command.CreateTodoCommand{Title: "Test", Description: "Desc", Priority: "high"}

	repo.On("Save", mock.AnythingOfType("*model.Todo")).Return(nil)

	id, err := uc.CreateTodoUseCase(cmd)
	assert.NotEmpty(t, id)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestCreateTodoUseCase_SaveError(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	cmd := command.CreateTodoCommand{Title: "Test", Description: "Desc", Priority: "high"}

	repo.On("Save", mock.AnythingOfType("*model.Todo")).Return(errors.New("db error"))

	id, err := uc.CreateTodoUseCase(cmd)
	assert.Empty(t, id)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to save todo", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestUpdateTodoUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todo := model.NewTodo("Original", "Desc", model.TodoPriorityMedium)
	cmd := command.UpdateTodoCommand{ID: "test-id", Title: "Updated"}

	repo.On("FindByID", model.TodoID("test-id")).Return(todo, nil)
	repo.On("Save", todo).Return(nil)

	err := uc.UpdateTodoUseCase(cmd)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateTodoUseCase_NotFound(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	cmd := command.UpdateTodoCommand{ID: "notfound", Title: "New Title"}

	repo.On("FindByID", model.TodoID("notfound")).Return(nil, errors.New("not found"))

	err := uc.UpdateTodoUseCase(cmd)
	assert.NotNil(t, err)
	assert.Equal(t, "Todo not found", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestUpdateTodoUseCase_InvalidTitle(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	// Create a title that's too long (over 100 characters - domain service limit)
	longTitle := "This is a very long title that exceeds the maximum allowed length of 100 characters. " +
		"It should trigger a validation error in the domain service."
	cmd := command.UpdateTodoCommand{ID: "test-id", Title: longTitle}

	// Note: FindByID is not called because domain validation fails first

	err := uc.UpdateTodoUseCase(cmd)
	assert.NotNil(t, err)
	assert.Equal(t, "Title too long", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestCompleteTodoUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todo := model.NewTodo("Test", "Desc", model.TodoPriorityMedium)

	repo.On("FindByID", todo.GetID()).Return(todo, nil)
	repo.On("Save", todo).Return(nil)

	err := uc.CompleteTodoUseCase(todo.GetID())
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestCompleteTodoUseCase_NotFound(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	id := model.TodoID("notfound")

	repo.On("FindByID", id).Return(nil, errors.New("not found"))

	err := uc.CompleteTodoUseCase(id)
	assert.NotNil(t, err)
	assert.Equal(t, "Todo not found", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestCompleteTodoUseCase_AlreadyCompleted(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todo := model.NewTodo("Already Done", "Desc", model.TodoPriorityMedium)
	todo.MarkAsCompleted() // Mark as completed first

	repo.On("FindByID", todo.GetID()).Return(todo, nil)

	err := uc.CompleteTodoUseCase(todo.GetID())
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot complete todo", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestArchiveTodoUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todo := model.NewTodo("Test", "Desc", model.TodoPriorityMedium)

	repo.On("FindByID", todo.GetID()).Return(todo, nil)
	repo.On("Save", todo).Return(nil)

	err := uc.ArchiveTodoUseCase(todo.GetID())
	assert.Nil(t, err)
	repo.AssertExpectations(t)
}

func TestArchiveTodoUseCase_NotFound(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	id := model.TodoID("notfound")

	repo.On("FindByID", id).Return(nil, errors.New("not found"))

	err := uc.ArchiveTodoUseCase(id)
	assert.NotNil(t, err)
	assert.Equal(t, "Todo not found", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestGetTodoUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todo := model.NewTodo("Test", "Desc", model.TodoPriorityMedium)
	repo.On("FindByID", todo.GetID()).Return(todo, nil)

	resp, err := uc.GetTodoUseCase(todo.GetID())
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, string(todo.GetID()), resp.ID)
	repo.AssertExpectations(t)
}

func TestGetTodoUseCase_NotFound(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	id := model.TodoID("notfound")
	repo.On("FindByID", id).Return(nil, errors.New("not found"))

	resp, err := uc.GetTodoUseCase(id)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "Todo not found", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestListTodosUseCase_Success(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	todos := []*model.Todo{
		model.NewTodo("Todo 1", "Desc 1", model.TodoPriorityHigh),
		model.NewTodo("Todo 2", "Desc 2", model.TodoPriorityMedium),
	}
	repo.On("FindAll").Return(todos, nil)

	resp, err := uc.ListTodosUseCase()
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.Equal(t, 2, resp.Count)
	repo.AssertExpectations(t)
}

func TestListTodosUseCase_RepoError(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)
	repo.On("FindAll").Return(nil, errors.New("db error"))

	resp, err := uc.ListTodosUseCase()
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to retrieve todos", err.GetErrorMessage())
	repo.AssertExpectations(t)
}

func TestTestErrorUseCase(t *testing.T) {
	repo := new(MockTodoRepository)
	domainService := service.NewTodoDomainService()
	uc := NewTodoUseCase(repo, domainService)

	err := uc.TestErrorUseCase()
	assert.NotNil(t, err)
	assert.Equal(t, "Test error message", err.GetErrorMessage())
	assert.Equal(t, 400, err.GetHttpStatus())
}
