package command

// CreateTodoCommand represents a command to create a new Todo following CQRS pattern
type CreateTodoCommand struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Priority    string `json:"priority,omitempty" validate:"oneof=low medium high"`
	CategoryID  string `json:"category-id,omitempty"`
	CreatedBy   string `json:"created-by,omitempty"`
}

// UpdateTodoCommand represents a command to update an existing Todo
type UpdateTodoCommand struct {
	ID          string `json:"id" validate:"required"`
	Title       string `json:"title,omitempty" validate:"max=200"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Priority    string `json:"priority,omitempty" validate:"oneof=low medium high"`
	CategoryID  string `json:"category-id,omitempty"`
}

// CompleteTodoCommand represents a command to mark a Todo as completed
type CompleteTodoCommand struct {
	ID string `json:"id" validate:"required"`
}

// ArchiveTodoCommand represents a command to archive a Todo
type ArchiveTodoCommand struct {
	ID string `json:"id" validate:"required"`
}

// CreateUserCommand represents a command to create a new User
type CreateUserCommand struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	FirstName string `json:"first-name" validate:"required,min=1,max=50"`
	LastName  string `json:"last-name" validate:"required,min=1,max=50"`
}

// UpdateUserProfileCommand represents a command to update user profile
type UpdateUserProfileCommand struct {
	ID        string `json:"id" validate:"required"`
	FirstName string `json:"first-name,omitempty" validate:"min=1,max=50"`
	LastName  string `json:"last-name,omitempty" validate:"min=1,max=50"`
	Email     string `json:"email,omitempty" validate:"email"`
}

// PromoteUserCommand represents a command to promote a user to admin
type PromoteUserCommand struct {
	ID string `json:"id" validate:"required"`
}

// SuspendUserCommand represents a command to suspend a user account
type SuspendUserCommand struct {
	ID string `json:"id" validate:"required"`
}

// CreateCategoryCommand represents a command to create a new Category
type CreateCategoryCommand struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Description string `json:"description,omitempty" validate:"max=200"`
	Color       string `json:"color" validate:"required,oneof=red blue green yellow purple orange gray"`
	CreatedBy   string `json:"created-by,omitempty"`
}

// UpdateCategoryCommand represents a command to update a Category
type UpdateCategoryCommand struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name,omitempty" validate:"min=1,max=50"`
	Description string `json:"description,omitempty" validate:"max=200"`
	Color       string `json:"color,omitempty" validate:"oneof=red blue green yellow purple orange gray"`
}
