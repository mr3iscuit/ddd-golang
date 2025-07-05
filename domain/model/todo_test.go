package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	todo := NewTodo("Test Title", "Test Description", TodoPriorityHigh)
	assert.NotEmpty(t, todo.GetID())
	assert.Equal(t, "Test Title", todo.GetTitle())
	assert.Equal(t, "Test Description", todo.GetDescription())
	assert.Equal(t, TodoStatusPending, todo.GetStatus())
	assert.Equal(t, TodoPriorityHigh, todo.GetPriority())
	assert.WithinDuration(t, time.Now(), todo.GetCreatedAt(), time.Second)
	assert.WithinDuration(t, time.Now(), todo.GetUpdatedAt(), time.Second)
	assert.Nil(t, todo.GetCompletedAt())
}

func TestMarkAsCompleted(t *testing.T) {
	todo := NewSimpleTodo("Complete Me")
	err := todo.MarkAsCompleted()
	assert.NoError(t, err)
	assert.Equal(t, TodoStatusCompleted, todo.GetStatus())
	assert.NotNil(t, todo.GetCompletedAt())

	// Marking again should fail
	err = todo.MarkAsCompleted()
	assert.Error(t, err)
}

func TestUpdateTitle(t *testing.T) {
	todo := NewSimpleTodo("Old Title")
	err := todo.UpdateTitle("New Title")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", todo.GetTitle())

	err = todo.UpdateTitle("")
	assert.Error(t, err)
}

func TestUpdateDescription(t *testing.T) {
	todo := NewSimpleTodo("Desc Test")
	err := todo.UpdateDescription("New Description")
	assert.NoError(t, err)
	assert.Equal(t, "New Description", todo.GetDescription())

	longDesc := make([]byte, 1001)
	for i := range longDesc {
		longDesc[i] = 'a'
	}
	err = todo.UpdateDescription(string(longDesc))
	assert.Error(t, err)
}

func TestUpdatePriority(t *testing.T) {
	todo := NewSimpleTodo("Priority Test")
	err := todo.UpdatePriority(TodoPriorityHigh)
	assert.NoError(t, err)
	assert.Equal(t, TodoPriorityHigh, todo.GetPriority())

	err = todo.UpdatePriority("invalid")
	assert.Error(t, err)
}

func TestArchiveTodo(t *testing.T) {
	todo := NewSimpleTodo("Archive Me")
	err := todo.ArchiveTodo()
	assert.NoError(t, err)
	assert.Equal(t, TodoStatusArchived, todo.GetStatus())

	err = todo.ArchiveTodo()
	assert.Error(t, err)
}
