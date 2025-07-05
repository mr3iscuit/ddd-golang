package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	handler "github.com/mr3iscuit/ddd-golang/adapters/http"
	"github.com/mr3iscuit/ddd-golang/application/usecase"
	"github.com/mr3iscuit/ddd-golang/domain/service"
	"github.com/mr3iscuit/ddd-golang/infrastructure/repository"
)

// Helper to start a test server and return its URL and a cleanup function
func startTestServer() (string, func()) {
	repo := repository.NewInMemoryTodoRepository()
	domainService := service.NewTodoDomainService()
	useCase := usecase.NewTodoUseCase(repo, domainService)
	h := handler.NewTodoHTTPAdapter(useCase)

	r := chi.NewRouter()
	r.Post("/todos", h.HandleCreateTodo)
	r.Get("/todos", h.HandleListTodos)
	r.Get("/todos/{id}", h.HandleGetTodo)
	r.Put("/todos/{id}", h.HandleUpdateTodo)
	r.Put("/todos/{id}/complete", h.HandleCompleteTodo)
	r.Put("/todos/{id}/archive", h.HandleArchiveTodo)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: r}
	go server.Serve(ln)
	url := "http://" + ln.Addr().String()
	cleanup := func() {
		server.Close()
		ln.Close()
	}
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	return url, cleanup
}

func TestIntegration_FullFlow_Curl(t *testing.T) {
	url, cleanup := startTestServer()
	defer cleanup()

	// 1. Create a todo
	createBody := `{"title":"Integration Test Todo","description":"This is a test todo","priority":"high"}`
	cmd := exec.Command("curl", "-s", "-X", "POST", "-H", "Content-Type: application/json", "-d", createBody, url+"/todos")
	out, err := cmd.Output()
	assert.NoError(t, err)

	var createResp map[string]string
	json.Unmarshal(out, &createResp)
	todoID := createResp["id"]
	assert.NotEmpty(t, todoID)

	// 2. List todos
	cmd = exec.Command("curl", "-s", url+"/todos")
	out, err = cmd.Output()
	assert.NoError(t, err)
	var listResp struct {
		Todos []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Priority string `json:"priority"`
		} `json:"todos"`
		Count int `json:"count"`
	}
	json.Unmarshal(out, &listResp)
	assert.Equal(t, 1, listResp.Count)
	assert.Equal(t, "Integration Test Todo", listResp.Todos[0].Title)

	// 3. Update the todo
	updateBody := `{"id":"` + todoID + `","title":"Updated Title","description":"Updated Description","priority":"high"}`
	cmd = exec.Command("curl", "-s", "-X", "PUT", "-H", "Content-Type: application/json", "-d", updateBody, url+"/todos/"+todoID)
	out, err = cmd.Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Todo updated successfully")

	// 4. Get the todo
	cmd = exec.Command("curl", "-s", url+"/todos/"+todoID)
	out, err = cmd.Output()
	assert.NoError(t, err)
	var getResp struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
	}
	json.Unmarshal(out, &getResp)
	assert.Equal(t, "Updated Title", getResp.Title)
	assert.Equal(t, "Updated Description", getResp.Description)
	assert.Equal(t, "high", getResp.Priority)
	assert.Equal(t, "pending", getResp.Status)

	// 5. Complete the todo
	cmd = exec.Command("curl", "-s", "-X", "PUT", url+"/todos/"+todoID+"/complete")
	out, err = cmd.Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Todo completed successfully")

	// 6. Get the completed todo
	cmd = exec.Command("curl", "-s", url+"/todos/"+todoID)
	out, err = cmd.Output()
	assert.NoError(t, err)
	var getCompletedResp struct {
		Status      string `json:"status"`
		CompletedAt string `json:"completed-at"`
	}
	json.Unmarshal(out, &getCompletedResp)
	assert.Equal(t, "completed", getCompletedResp.Status)
	assert.NotEmpty(t, getCompletedResp.CompletedAt)

	// 7. Archive the todo
	cmd = exec.Command("curl", "-s", "-X", "PUT", url+"/todos/"+todoID+"/archive")
	out, err = cmd.Output()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "Todo archived successfully")

	// 8. Get the archived todo
	cmd = exec.Command("curl", "-s", url+"/todos/"+todoID)
	out, err = cmd.Output()
	assert.NoError(t, err)
	var getArchivedResp struct {
		Status string `json:"status"`
	}
	json.Unmarshal(out, &getArchivedResp)
	assert.Equal(t, "archived", getArchivedResp.Status)
}
