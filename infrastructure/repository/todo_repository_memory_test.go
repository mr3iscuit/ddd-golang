package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mr3iscuit/ddd-golang/domain/model"
)

func TestNewInMemoryTodoRepository(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.todos)
	assert.Equal(t, 0, len(repo.todos))
}

func TestSaveAndFindByID(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	todo := model.NewTodo("Test Todo", "Test Description", model.TodoPriorityHigh)

	// Save todo
	err := repo.Save(todo)
	assert.NoError(t, err)

	// Find by ID
	found, err := repo.FindByID(todo.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, todo.GetID(), found.GetID())
	assert.Equal(t, todo.GetTitle(), found.GetTitle())
}

func TestFindByID_NotFound(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	id := model.TodoID("non-existent-id")

	found, err := repo.FindByID(id)
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "not found")
}

func TestFindAll(t *testing.T) {
	repo := NewInMemoryTodoRepository()

	// Add multiple todos
	todo1 := model.NewTodo("Todo 1", "Desc 1", model.TodoPriorityLow)
	todo2 := model.NewTodo("Todo 2", "Desc 2", model.TodoPriorityMedium)
	todo3 := model.NewTodo("Todo 3", "Desc 3", model.TodoPriorityHigh)

	repo.Save(todo1)
	repo.Save(todo2)
	repo.Save(todo3)

	// Find all
	todos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(todos))
}

func TestFindAll_Empty(t *testing.T) {
	repo := NewInMemoryTodoRepository()

	todos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(todos))
}

func TestDelete(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	todo := model.NewTodo("To Delete", "Will be deleted", model.TodoPriorityMedium)

	// Save first
	err := repo.Save(todo)
	assert.NoError(t, err)

	// Verify it exists
	found, err := repo.FindByID(todo.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, found)

	// Delete
	err = repo.Delete(todo.GetID())
	assert.NoError(t, err)

	// Verify it's gone
	found, err = repo.FindByID(todo.GetID())
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestDelete_NotFound(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	id := model.TodoID("non-existent-id")

	err := repo.Delete(id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestConcurrentAccess(t *testing.T) {
	repo := NewInMemoryTodoRepository()

	// Test concurrent saves
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			todo := model.NewTodo("Concurrent Todo", "Description", model.TodoPriorityMedium)
			err := repo.Save(todo)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all todos were saved
	todos, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(todos))
}

func TestUpdateExistingTodo(t *testing.T) {
	repo := NewInMemoryTodoRepository()
	todo := model.NewTodo("Original Title", "Original Description", model.TodoPriorityLow)

	// Save original
	err := repo.Save(todo)
	assert.NoError(t, err)

	// Update the todo
	todo.UpdateTitle("Updated Title")
	todo.UpdateDescription("Updated Description")
	todo.UpdatePriority(model.TodoPriorityHigh)

	// Save updated version
	err = repo.Save(todo)
	assert.NoError(t, err)

	// Retrieve and verify updates
	found, err := repo.FindByID(todo.GetID())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", found.GetTitle())
	assert.Equal(t, "Updated Description", found.GetDescription())
	assert.Equal(t, model.TodoPriorityHigh, found.GetPriority())
}
