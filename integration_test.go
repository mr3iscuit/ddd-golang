package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	handler "github.com/mr3iscuit/ddd-golang/adapters/http"
	"github.com/mr3iscuit/ddd-golang/application/usecase"
	"github.com/mr3iscuit/ddd-golang/domain/service"
	postgresrepo "github.com/mr3iscuit/ddd-golang/infrastructure/repository/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mr3iscuit/ddd-golang/pkg/config"
)

func startPostgresTestServer(t *testing.T) (string, func()) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading configuration for tests: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to Postgres with GORM: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&postgresrepo.TodoRecord{})
	if err != nil {
		t.Fatalf("Failed to auto-migrate schema: %v", err)
	}

	// Clean up all todos before test
	if err := db.Exec("DELETE FROM todos").Error; err != nil {
		t.Fatalf("Failed to clean todos table: %v", err)
	}

	repo := postgresrepo.NewPostgresTodoRepository(db)
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
		t.Fatal(err)
	}
	server := &http.Server{Handler: r}
	go server.Serve(ln)
	url := "http://" + ln.Addr().String()
	cleanup := func() {
		server.Close()
		ln.Close()

		// Drop the table after tests
		if err := db.Migrator().DropTable(&postgresrepo.TodoRecord{}); err != nil {
			t.Logf("Failed to drop table in cleanup: %v", err)
		}
	}
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	return url, cleanup
}

func TestIntegration_Postgres_FullFlow(t *testing.T) {
	url, cleanup := startPostgresTestServer(t)
	defer cleanup()

	client := &http.Client{}

	// 1. Create a todo
	createBody := map[string]string{"title": "PG Integration Test Todo", "description": "This is a test todo", "priority": "high"}
	createBodyBytes, _ := json.Marshal(createBody)
	req, _ := http.NewRequest(http.MethodPost, url+"/todos", bytes.NewBuffer(createBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResp map[string]string
	json.NewDecoder(resp.Body).Decode(&createResp)
	todoID := createResp["id"]
	assert.NotEmpty(t, todoID)
	resp.Body.Close()

	// 2. List todos
	req, _ = http.NewRequest(http.MethodGet, url+"/todos", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var listResp struct {
		Todos []struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Priority string `json:"priority"`
		} `json:"todos"`
		Count int `json:"count"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.GreaterOrEqual(t, listResp.Count, 1)
	found := false
	for _, todo := range listResp.Todos {
		if todo.ID == todoID {
			found = true
			assert.Equal(t, "PG Integration Test Todo", todo.Title)
		}
	}
	assert.True(t, found)
	resp.Body.Close()

	// 3. Update the todo
	updateBody := map[string]string{"id": todoID, "title": "Updated PG Title", "description": "Updated PG Description", "priority": "high"}
	updateBodyBytes, _ := json.Marshal(updateBody)
	req, _ = http.NewRequest(http.MethodPut, url+"/todos/"+todoID, bytes.NewBuffer(updateBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var updateResp map[string]string
	json.NewDecoder(resp.Body).Decode(&updateResp)
	assert.Equal(t, "Todo updated successfully", updateResp["message"])
	resp.Body.Close()

	// 4. Get the todo
	req, _ = http.NewRequest(http.MethodGet, url+"/todos/"+todoID, nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var getResp struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Status      string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&getResp)
	assert.Equal(t, "Updated PG Title", getResp.Title)
	assert.Equal(t, "Updated PG Description", getResp.Description)
	assert.Equal(t, "high", getResp.Priority)
	assert.Equal(t, "pending", getResp.Status)
	resp.Body.Close()

	// 5. Complete the todo
	req, _ = http.NewRequest(http.MethodPut, url+"/todos/"+todoID+"/complete", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var completeResp map[string]string
	json.NewDecoder(resp.Body).Decode(&completeResp)
	assert.Equal(t, "Todo completed successfully", completeResp["message"])
	resp.Body.Close()

	// 6. Get the completed todo
	req, _ = http.NewRequest(http.MethodGet, url+"/todos/"+todoID, nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var getCompletedResp struct {
		Status      string `json:"status"`
		CompletedAt string `json:"completed-at"`
	}
	json.NewDecoder(resp.Body).Decode(&getCompletedResp)
	assert.Equal(t, "completed", getCompletedResp.Status)
	assert.NotEmpty(t, getCompletedResp.CompletedAt)
	resp.Body.Close()

	// 7. Archive the todo
	req, _ = http.NewRequest(http.MethodPut, url+"/todos/"+todoID+"/archive", nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var archiveResp map[string]string
	json.NewDecoder(resp.Body).Decode(&archiveResp)
	assert.Equal(t, "Todo archived successfully", archiveResp["message"])
	resp.Body.Close()

	// 8. Get the archived todo
	req, _ = http.NewRequest(http.MethodGet, url+"/todos/"+todoID, nil)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var getArchivedResp struct {
		Status string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&getArchivedResp)
	assert.Equal(t, "archived", getArchivedResp.Status)
	resp.Body.Close()
}
