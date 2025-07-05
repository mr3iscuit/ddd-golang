package usecase

import (
	"github.com/mr3iscuit/ddd-golang/application/command"
	appmodel "github.com/mr3iscuit/ddd-golang/application/model"
	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
	"github.com/mr3iscuit/ddd-golang/domain/service"
)

// TodoUseCase implements the TodoUseCasePort
// and uses the TodoRepositoryPort
// (was TodoApplicationService)
type TodoUseCase struct {
	todoRepo      port.TodoRepositoryPort
	domainService *service.TodoDomainService
}

func NewTodoUseCase(todoRepo port.TodoRepositoryPort) *TodoUseCase {
	return &TodoUseCase{
		todoRepo:      todoRepo,
		domainService: service.NewTodoDomainService(),
	}
}

func (uc *TodoUseCase) CreateTodoUseCase(cmd command.CreateTodoCommand) (model.TodoID, *model.DomainError) {
	// Validate using domain service
	if err := uc.domainService.ValidateCreateTodoCommand(cmd.Title, cmd.Description, cmd.Priority); err != nil {
		return "", err
	}

	// Map priority string to domain type
	var priority model.TodoPriority
	switch cmd.Priority {
	case "low":
		priority = model.TodoPriorityLow
	case "high":
		priority = model.TodoPriorityHigh
	default:
		priority = model.TodoPriorityMedium
	}

	todo := model.NewTodo(cmd.Title, cmd.Description, priority)
	if err := uc.todoRepo.Save(todo); err != nil {
		return "", model.ErrFailedToSaveTodo
	}
	return todo.GetID(), nil
}

func (uc *TodoUseCase) UpdateTodoUseCase(cmd command.UpdateTodoCommand) *model.DomainError {
	// Validate using domain service
	if err := uc.domainService.ValidateUpdateTodoCommand(cmd.Title, cmd.Description, cmd.Priority); err != nil {
		return err
	}

	todo, err := uc.todoRepo.FindByID(model.TodoID(cmd.ID))
	if err != nil {
		return model.ErrTodoNotFound
	}

	if cmd.Title != "" {
		if err := todo.UpdateTitle(cmd.Title); err != nil {
			return model.ErrInvalidTitle
		}
	}

	if cmd.Description != "" {
		if err := todo.UpdateDescription(cmd.Description); err != nil {
			return model.ErrInvalidDescription
		}
	}

	if cmd.Priority != "" {
		var priority model.TodoPriority
		switch cmd.Priority {
		case "low":
			priority = model.TodoPriorityLow
		case "high":
			priority = model.TodoPriorityHigh
		case "medium":
			priority = model.TodoPriorityMedium
		default:
			return model.ErrInvalidPriority
		}
		if err := todo.UpdatePriority(priority); err != nil {
			return model.ErrInvalidPriority
		}
	}

	if err := uc.todoRepo.Save(todo); err != nil {
		return model.ErrFailedToSaveTodo
	}
	return nil
}

func (uc *TodoUseCase) CompleteTodoUseCase(id model.TodoID) *model.DomainError {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return model.ErrTodoNotFound
	}
	if err := todo.MarkAsCompleted(); err != nil {
		return model.ErrCannotCompleteTodo
	}
	if err := uc.todoRepo.Save(todo); err != nil {
		return model.ErrFailedToSaveCompletedTodo
	}
	return nil
}

func (uc *TodoUseCase) ArchiveTodoUseCase(id model.TodoID) *model.DomainError {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return model.ErrTodoNotFound
	}
	if err := todo.ArchiveTodo(); err != nil {
		return model.ErrCannotArchiveTodo
	}
	if err := uc.todoRepo.Save(todo); err != nil {
		return model.ErrFailedToSaveArchivedTodo
	}
	return nil
}

func (uc *TodoUseCase) GetTodoUseCase(id model.TodoID) (*appmodel.TodoResponse, *model.DomainError) {
	todo, err := uc.todoRepo.FindByID(id)
	if err != nil {
		return nil, model.ErrTodoNotFound
	}
	response := appmodel.TodoResponseMapper(todo)
	return &response, nil
}

func (uc *TodoUseCase) ListTodosUseCase() (*appmodel.TodoListResponse, *model.DomainError) {
	if uc.todoRepo == nil {
		return nil, model.ErrRepositoryNotInitialized
	}
	todos, err := uc.todoRepo.FindAll()
	if err != nil {
		return nil, model.ErrFailedToRetrieveTodos
	}
	response := appmodel.TodoListResponseMapper(todos)
	return &response, nil
}

func (uc *TodoUseCase) TestErrorUseCase() *model.DomainError {
	return model.ErrTestError
}
