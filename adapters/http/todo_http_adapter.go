package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mr3iscuit/ddd-golang/application/command"
	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/mr3iscuit/ddd-golang/pkg/config"
)

// TodoHTTPAdapter implements HTTP endpoints using the TodoUseCasePort
type TodoHTTPAdapter struct {
	usecase port.TodoUseCasePort
	config  *config.Config
}

// NewTodoHTTPAdapter creates a new Todo HTTP handler
func NewTodoHTTPAdapter(usecase port.TodoUseCasePort, cfg *config.Config) *TodoHTTPAdapter {
	return &TodoHTTPAdapter{usecase: usecase, config: cfg}
}

// writeJSONResponse writes a JSON response with the given status code
func (h *TodoHTTPAdapter) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeDomainError writes a domain error as JSON response
func (h *TodoHTTPAdapter) writeDomainError(w http.ResponseWriter, err model.DomainErrorPort) {
	errorResponse := err.ToResponse()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Error-Type", "domain-error")
	w.WriteHeader(err.GetHttpStatus())
	json.NewEncoder(w).Encode(errorResponse)
}

// parseJSON parses JSON from request body
func (h *TodoHTTPAdapter) parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *TodoHTTPAdapter) Router() http.Handler {
	r := chi.NewRouter()

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", h.config.ServerPort)),
	))

	// Todo endpoints
	r.Get("/todos", h.HandleListTodos)
	r.Post("/todos", h.HandleCreateTodo)
	r.Get("/todos/{id}", h.HandleGetTodo)
	r.Put("/todos/{id}", h.HandleUpdateTodo)
	r.Put("/todos/{id}/complete", h.HandleCompleteTodo)
	r.Put("/todos/{id}/archive", h.HandleArchiveTodo)

	// Test endpoint that always returns an error
	r.Get("/test-error", h.HandleTestError)
	return r
}

// HandleListTodos handles GET /todos
// @Summary List all todos
// @Description Get all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} appmodel.TodoResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos [get]
func (h *TodoHTTPAdapter) HandleListTodos(w http.ResponseWriter, r *http.Request) {
	response, err := h.usecase.ListTodosUseCase()
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// HandleCreateTodo handles POST /todos
// @Summary Create a new todo
// @Description Create a new todo with the given details
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body command.CreateTodoCommand true "Todo to create"
// @Success 201 {object} map[string]string
// @Failure 400 {object} appmodel.ErrorResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos [post]
func (h *TodoHTTPAdapter) HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var cmd command.CreateTodoCommand
	if err := h.parseJSON(r, &cmd); err != nil {
		h.writeDomainError(w, model.ErrInvalidJSON)
		return
	}

	id, err := h.usecase.CreateTodoUseCase(cmd)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, map[string]string{"id": string(id)})
}

// HandleGetTodo handles GET /todos/{id}
// @Summary Get a todo by ID
// @Description Get a specific todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} appmodel.TodoResponse
// @Failure 404 {object} appmodel.ErrorResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHTTPAdapter) HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeDomainError(w, model.ErrTodoNotFound)
		return
	}

	response, err := h.usecase.GetTodoUseCase(model.TodoID(id))
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// HandleUpdateTodo handles PUT /todos/{id}
// @Summary Update a todo
// @Description Update an existing todo
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body command.UpdateTodoCommand true "Todo updates"
// @Success 200 {object} map[string]string
// @Failure 400 {object} appmodel.ErrorResponse
// @Failure 404 {object} appmodel.ErrorResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos/{id} [put]
func (h *TodoHTTPAdapter) HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeDomainError(w, model.ErrTodoNotFound)
		return
	}

	var cmd command.UpdateTodoCommand
	if err := h.parseJSON(r, &cmd); err != nil {
		h.writeDomainError(w, model.ErrInvalidJSON)
		return
	}

	cmd.ID = id
	err := h.usecase.UpdateTodoUseCase(cmd)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Todo updated successfully"})
}

// HandleCompleteTodo handles PUT /todos/{id}/complete
// @Summary Complete a todo
// @Description Mark a todo as completed
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} appmodel.ErrorResponse
// @Failure 404 {object} appmodel.ErrorResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos/{id}/complete [put]
func (h *TodoHTTPAdapter) HandleCompleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeDomainError(w, model.ErrTodoNotFound)
		return
	}

	err := h.usecase.CompleteTodoUseCase(model.TodoID(id))
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Todo completed successfully"})
}

// HandleArchiveTodo handles PUT /todos/{id}/archive
// @Summary Archive a todo
// @Description Mark a todo as archived
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} appmodel.ErrorResponse
// @Failure 404 {object} appmodel.ErrorResponse
// @Failure 500 {object} appmodel.ErrorResponse
// @Router /todos/{id}/archive [put]
func (h *TodoHTTPAdapter) HandleArchiveTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.writeDomainError(w, model.ErrTodoNotFound)
		return
	}

	err := h.usecase.ArchiveTodoUseCase(model.TodoID(id))
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Todo archived successfully"})
}

// HandleTestError handles GET /test-error
// @Summary Test error endpoint
// @Description Returns a test error for testing error handling
// @Tags test
// @Produce json
// @Success 400 {object} appmodel.ErrorResponse
// @Router /test-error [get]
func (h *TodoHTTPAdapter) HandleTestError(w http.ResponseWriter, r *http.Request) {
	err := h.usecase.TestErrorUseCase()
	h.writeDomainError(w, err)
}
