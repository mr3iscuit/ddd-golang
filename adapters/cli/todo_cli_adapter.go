package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mr3iscuit/ddd-golang/application/command"
	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// TodoCLIAdapter handles command-line interface for Todo operations
type TodoCLIAdapter struct {
	usecase port.TodoUseCasePort
}

// NewTodoCLIAdapter creates a new Todo CLI
func NewTodoCLIAdapter(usecase port.TodoUseCasePort) *TodoCLIAdapter {
	return &TodoCLIAdapter{usecase: usecase}
}

// Run starts the CLI application
func (c *TodoCLIAdapter) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Todo CLI - Type 'help' for commands")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" || input == "exit" {
			break
		}

		c.handleCommand(input)
	}
}

// handleCommand processes user input commands
func (c *TodoCLIAdapter) handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "add":
		if len(parts) < 2 {
			fmt.Println("Usage: add <title> [description] [priority]")
			return
		}
		title := parts[1]
		description := ""
		priority := "medium"

		if len(parts) > 2 {
			description = parts[2]
		}
		if len(parts) > 3 {
			priority = parts[3]
		}

		cmd := command.CreateTodoCommand{
			Title:       title,
			Description: description,
			Priority:    priority,
		}
		id, err := c.usecase.CreateTodoUseCase(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
		} else {
			fmt.Printf("Todo created with ID: %s\n", id)
		}

	case "list":
		todoListResponse, err := c.usecase.ListTodosUseCase()
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
			return
		}
		if todoListResponse.Count == 0 {
			fmt.Println("No todos found")
			return
		}
		fmt.Printf("Found %d todos:\n", todoListResponse.Count)
		for _, todo := range todoListResponse.Todos {
			status := todo.Status
			priority := todo.Priority
			fmt.Printf("[%s] %s - %s (Priority: %s)\n", todo.ID, todo.Title, status, priority)
		}

	case "get":
		if len(parts) < 2 {
			fmt.Println("Usage: get <id>")
			return
		}
		todoID := model.TodoID(parts[1])
		todoResponse, err := c.usecase.GetTodoUseCase(todoID)
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
			return
		}
		fmt.Printf("Todo Details:\n")
		fmt.Printf("  ID: %s\n", todoResponse.ID)
		fmt.Printf("  Title: %s\n", todoResponse.Title)
		fmt.Printf("  Description: %s\n", todoResponse.Description)
		fmt.Printf("  Status: %s\n", todoResponse.Status)
		fmt.Printf("  Priority: %s\n", todoResponse.Priority)
		fmt.Printf("  Created: %s\n", todoResponse.CreatedAt.Format("2006-01-02 15:04:05"))
		if todoResponse.CompletedAt != nil {
			fmt.Printf("  Completed: %s\n", todoResponse.CompletedAt.Format("2006-01-02 15:04:05"))
		}

	case "update":
		if len(parts) < 3 {
			fmt.Println("Usage: update <id> <title> [description] [priority]")
			return
		}
		id := parts[1]
		title := parts[2]
		description := ""
		priority := ""

		if len(parts) > 3 {
			description = parts[3]
		}
		if len(parts) > 4 {
			priority = parts[4]
		}

		cmd := command.UpdateTodoCommand{
			ID:          id,
			Title:       title,
			Description: description,
			Priority:    priority,
		}
		err := c.usecase.UpdateTodoUseCase(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
		} else {
			fmt.Println("Todo updated successfully")
		}

	case "complete":
		if len(parts) < 2 {
			fmt.Println("Usage: complete <id>")
			return
		}
		err := c.usecase.CompleteTodoUseCase(model.TodoID(parts[1]))
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
		} else {
			fmt.Println("Todo completed successfully")
		}

	case "archive":
		if len(parts) < 2 {
			fmt.Println("Usage: archive <id>")
			return
		}
		err := c.usecase.ArchiveTodoUseCase(model.TodoID(parts[1]))
		if err != nil {
			fmt.Printf("Error: %s\n", err.GetErrorMessage())
		} else {
			fmt.Println("Todo archived successfully")
		}

	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  add <title> [description] [priority] - Add a new todo")
		fmt.Println("  list                                - List all todos")
		fmt.Println("  get <id>                           - Get todo details")
		fmt.Println("  update <id> <title> [desc] [priority] - Update a todo")
		fmt.Println("  complete <id>                      - Complete a todo")
		fmt.Println("  archive <id>                       - Archive a todo")
		fmt.Println("  help                               - Show this help")
		fmt.Println("  quit/exit                          - Exit the application")
		fmt.Println("\nPriority options: low, medium, high")

	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", parts[0])
	}
}
