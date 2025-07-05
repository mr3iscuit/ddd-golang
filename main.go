// @title           Todo API
// @version         1.0
// @description     A DDD-style Todo API with proper layering
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
// @schemes http

// @securityDefinitions.basic  BasicAuth
package main

import (
	"log"
	"net/http"

	handler "github.com/mr3iscuit/ddd-golang/adapters/http"
	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/application/usecase"
	_ "github.com/mr3iscuit/ddd-golang/docs"
	"github.com/mr3iscuit/ddd-golang/infrastructure/repository"
)

func main() {
	// Outbound port (repository)
	var todoRepo port.TodoRepositoryPort = repository.NewInMemoryTodoRepository()
	// Use case (inbound port implementation)
	var todoUseCase port.TodoUseCasePort = usecase.NewTodoUseCase(todoRepo)
	// Handler (inbound adapter)
	todoHandler := handler.NewTodoHTTPAdapter(todoUseCase)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", todoHandler.Router()); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// CLI usage (uncomment to use CLI instead of HTTP)
	// cli := cli.NewCLI(todoService)
	// cli.Run()
}
