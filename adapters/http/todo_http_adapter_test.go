package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mr3iscuit/ddd-golang/application/command"
	appmodel "github.com/mr3iscuit/ddd-golang/application/model"
	"github.com/mr3iscuit/ddd-golang/domain/model"
	"github.com/mr3iscuit/ddd-golang/pkg/config"
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

func TestHandleCreateTodo_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	cmd := command.CreateTodoCommand{
		Title:       "Test Todo",
		Description: "Test Description",
		Priority:    "high",
	}

	mockUseCase.On("CreateTodoUseCase", cmd).Return(model.TodoID("test-id"), (*model.DomainError)(nil))

	body, _ := json.Marshal(cmd)
	req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreateTodo(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "test-id", response["id"])

	mockUseCase.AssertExpectations(t)
}

func TestHandleCreateTodo_InvalidJSON(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreateTodo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response appmodel.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid JSON", response.ErrorMessage)
}

func TestHandleCreateTodo_UseCaseError(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	cmd := command.CreateTodoCommand{Title: "Test"}
	domainError := model.NewDomainError(1001, 400, "Validation failed", "Title too short", nil)

	mockUseCase.On("CreateTodoUseCase", cmd).Return(model.TodoID(""), domainError)

	body, _ := json.Marshal(cmd)
	req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreateTodo(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response appmodel.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Validation failed", response.ErrorMessage)

	mockUseCase.AssertExpectations(t)
}

func TestHandleListTodos_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	todos := []appmodel.TodoResponse{
		{ID: "1", Title: "Todo 1", Status: "pending", Priority: "high"},
		{ID: "2", Title: "Todo 2", Status: "completed", Priority: "medium"},
	}
	response := &appmodel.TodoListResponse{Todos: todos, Count: 2}

	mockUseCase.On("ListTodosUseCase").Return(response, (*model.DomainError)(nil))

	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	handler.HandleListTodos(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result appmodel.TodoListResponse
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 2, result.Count)
	assert.Equal(t, "Todo 1", result.Todos[0].Title)

	mockUseCase.AssertExpectations(t)
}

func TestHandleListTodos_UseCaseError(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	domainError := model.NewDomainError(4001, 500, "Database error", "Connection failed", nil)
	mockUseCase.On("ListTodosUseCase").Return((*appmodel.TodoListResponse)(nil), domainError)

	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	handler.HandleListTodos(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response appmodel.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Database error", response.ErrorMessage)

	mockUseCase.AssertExpectations(t)
}

func TestHandleGetTodo_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	todoID := model.TodoID("test-id")
	todoResponse := &appmodel.TodoResponse{
		ID:       "test-id",
		Title:    "Test Todo",
		Status:   "pending",
		Priority: "high",
	}

	mockUseCase.On("GetTodoUseCase", todoID).Return(todoResponse, (*model.DomainError)(nil))

	req := httptest.NewRequest("GET", "/todos/test-id", nil)
	w := httptest.NewRecorder()

	// Create a chi router to properly handle URL parameters
	r := chi.NewRouter()
	r.Get("/todos/{id}", handler.HandleGetTodo)

	// Serve the request through the router
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result appmodel.TodoResponse
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "test-id", result.ID)
	assert.Equal(t, "Test Todo", result.Title)

	mockUseCase.AssertExpectations(t)
}

func TestHandleCompleteTodo_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	todoID := model.TodoID("test-id")
	mockUseCase.On("CompleteTodoUseCase", todoID).Return((*model.DomainError)(nil))

	req := httptest.NewRequest("PUT", "/todos/test-id/complete", nil)
	w := httptest.NewRecorder()

	// Create a chi router to properly handle URL parameters
	r := chi.NewRouter()
	r.Put("/todos/{id}/complete", handler.HandleCompleteTodo)

	// Serve the request through the router
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Todo completed successfully", response["message"])

	mockUseCase.AssertExpectations(t)
}

func TestHandleArchiveTodo_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	todoID := model.TodoID("test-id")
	mockUseCase.On("ArchiveTodoUseCase", todoID).Return((*model.DomainError)(nil))

	req := httptest.NewRequest("PUT", "/todos/test-id/archive", nil)
	w := httptest.NewRecorder()

	// Create a chi router to properly handle URL parameters
	r := chi.NewRouter()
	r.Put("/todos/{id}/archive", handler.HandleArchiveTodo)

	// Serve the request through the router
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Todo archived successfully", response["message"])

	mockUseCase.AssertExpectations(t)
}

func TestHandleUpdateTodo_Success(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	cmd := command.UpdateTodoCommand{
		ID:          "test-id",
		Title:       "Updated Todo",
		Description: "Updated Description",
		Priority:    "medium",
	}

	mockUseCase.On("UpdateTodoUseCase", cmd).Return((*model.DomainError)(nil))

	body, _ := json.Marshal(cmd)
	req := httptest.NewRequest("PUT", "/todos/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create a chi router to properly handle URL parameters
	r := chi.NewRouter()
	r.Put("/todos/{id}", handler.HandleUpdateTodo)

	// Serve the request through the router
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Todo updated successfully", response["message"])

	mockUseCase.AssertExpectations(t)
}

func TestHandleTestError(t *testing.T) {
	mockUseCase := new(MockTodoUseCase)
	handler := NewTodoHTTPAdapter(mockUseCase, &config.Config{ServerPort: "8080"})

	domainError := model.NewDomainError(9001, 400, "Test error", "Test reason", nil)
	mockUseCase.On("TestErrorUseCase").Return(domainError)

	req := httptest.NewRequest("GET", "/test-error", nil)
	w := httptest.NewRecorder()

	handler.HandleTestError(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response appmodel.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Test error", response.ErrorMessage)

	mockUseCase.AssertExpectations(t)
}
