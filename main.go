// @title           Todo API
// @version         1.0
// @description     A DDD-style Todo API with proper layering
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /
// @schemes http

// @securityDefinitions.basic  BasicAuth
package main

import (
	"fmt"
	"log"
	"net/http"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	handler "github.com/mr3iscuit/ddd-golang/adapters/http"
	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/application/usecase"
	_ "github.com/mr3iscuit/ddd-golang/docs"
	"github.com/mr3iscuit/ddd-golang/domain/service"
	postgresrepo "github.com/mr3iscuit/ddd-golang/infrastructure/repository/postgres"

	"github.com/mr3iscuit/ddd-golang/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Outbound port (repository)
	var todoRepo port.TodoRepositoryPort

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Using PostgresTodoRepository")
	todoRepo = postgresrepo.NewPostgresTodoRepository(db)

	// Domain service (outbound port implementation)
	var domainService port.TodoDomainServicePort = service.NewTodoDomainService()
	// Use case (inbound port implementation)
	var todoUseCase port.TodoUseCasePort = usecase.NewTodoUseCase(todoRepo, domainService)
	// Handler (inbound adapter)
	todoHandler := handler.NewTodoHTTPAdapter(todoUseCase, cfg)

	log.Printf("Starting HTTP server on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), todoHandler.Router()); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// CLI usage (uncomment to use CLI instead of HTTP)
	// cli := cli.NewCLI(todoService)
	// cli.Run()
}
