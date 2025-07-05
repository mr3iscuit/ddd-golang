package cli

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/mr3iscuit/ddd-golang/application/command"
	appmodel "github.com/mr3iscuit/ddd-golang/application/model"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

type MockTodoUseCase struct {
	mock.Mock
}

func (m *MockTodoUseCase) CreateTodoUseCase(cmd command.CreateTodoCommand) (model.TodoID, *model.DomainError) {
	args := m.Called(cmd)
	return args.Get(0).(model.TodoID), args.Get(1).(*model.DomainError)
}

func (m *MockTodoUseCase) UpdateTodoUseCase(cmd command.UpdateTodoCommand) *model.DomainError {
	args := m.Called(cmd)
	return args.Get(0).(*model.DomainError)
}

func (m *MockTodoUseCase) CompleteTodoUseCase(id model.TodoID) *model.DomainError {
	args := m.Called(id)
	return args.Get(0).(*model.DomainError)
}

func (m *MockTodoUseCase) ArchiveTodoUseCase(id model.TodoID) *model.DomainError {
	args := m.Called(id)
	return args.Get(0).(*model.DomainError)
}

func (m *MockTodoUseCase) GetTodoUseCase(id model.TodoID) (*appmodel.TodoResponse, *model.DomainError) {
	args := m.Called(id)
	if resp, ok := args.Get(0).(*appmodel.TodoResponse); ok {
		return resp, args.Get(1).(*model.DomainError)
	}
	return nil, args.Get(1).(*model.DomainError)
}

func (m *MockTodoUseCase) ListTodosUseCase() (*appmodel.TodoListResponse, *model.DomainError) {
	args := m.Called()
	if resp, ok := args.Get(0).(*appmodel.TodoListResponse); ok {
		return resp, args.Get(1).(*model.DomainError)
	}
	return nil, args.Get(1).(*model.DomainError)
}

func (m *MockTodoUseCase) TestErrorUseCase() *model.DomainError {
	args := m.Called()
	return args.Get(0).(*model.DomainError)
}

func TestHandleCommand_Add(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	expectedCmd := command.CreateTodoCommand{
		Title:       "Test",
		Description: "Todo",
		Priority:    "Test",
	}

	mockUseCase.On("CreateTodoUseCase", expectedCmd).Return(model.TodoID("test-id"), (*model.DomainError)(nil))

	adapter.handleCommand("add Test Todo Test")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Add_Error(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	expectedCmd := command.CreateTodoCommand{
		Title:       "Test",
		Description: "Todo",
		Priority:    "medium",
	}

	domainError := model.NewDomainError(1001, 400, "Validation failed", "Title too short", nil)
	mockUseCase.On("CreateTodoUseCase", expectedCmd).Return(model.TodoID(""), domainError)

	adapter.handleCommand("add Test Todo")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_List_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	todos := []appmodel.TodoResponse{
		{ID: "1", Title: "Todo 1", Status: "pending", Priority: "high"},
		{ID: "2", Title: "Todo 2", Status: "completed", Priority: "medium"},
	}
	response := &appmodel.TodoListResponse{Todos: todos, Count: 2}

	mockUseCase.On("ListTodosUseCase").Return(response, (*model.DomainError)(nil))

	adapter.handleCommand("list")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_List_Empty(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	response := &appmodel.TodoListResponse{Todos: []appmodel.TodoResponse{}, Count: 0}
	mockUseCase.On("ListTodosUseCase").Return(response, (*model.DomainError)(nil))

	adapter.handleCommand("list")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Get_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	todoID := model.TodoID("test-id")
	todoResponse := &appmodel.TodoResponse{
		ID:          "test-id",
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
	}

	mockUseCase.On("GetTodoUseCase", todoID).Return(todoResponse, (*model.DomainError)(nil))

	adapter.handleCommand("get test-id")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Update_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	expectedCmd := command.UpdateTodoCommand{
		ID:          "test-id",
		Title:       "Updated",
		Description: "Title",
		Priority:    "Updated",
	}

	mockUseCase.On("UpdateTodoUseCase", expectedCmd).Return((*model.DomainError)(nil))

	adapter.handleCommand("update test-id Updated Title Updated")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Complete_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	todoID := model.TodoID("test-id")
	mockUseCase.On("CompleteTodoUseCase", todoID).Return((*model.DomainError)(nil))

	adapter.handleCommand("complete test-id")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Archive_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	todoID := model.TodoID("test-id")
	mockUseCase.On("ArchiveTodoUseCase", todoID).Return((*model.DomainError)(nil))

	adapter.handleCommand("archive test-id")

	mockUseCase.AssertExpectations(t)
}

func TestHandleCommand_Empty(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	// Should not call any use case methods
	adapter.handleCommand("")

	mockUseCase.AssertNotCalled(t, "CreateTodoUseCase")
	mockUseCase.AssertNotCalled(t, "ListTodosUseCase")
	mockUseCase.AssertNotCalled(t, "GetTodoUseCase")
}

func TestHandleCommand_Unknown(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	adapter := NewTodoCLIAdapter(mockUseCase)

	// Should not call any use case methods
	adapter.handleCommand("unknown")

	mockUseCase.AssertNotCalled(t, "CreateTodoUseCase")
	mockUseCase.AssertNotCalled(t, "ListTodosUseCase")
	mockUseCase.AssertNotCalled(t, "GetTodoUseCase")
}
