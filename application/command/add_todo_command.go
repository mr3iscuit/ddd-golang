package command

// CreateTodoCommand represents a command to create a new Todo following CQRS pattern
type CreateTodoCommand struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"`
	CategoryID  string `json:"category-id,omitempty"`
	CreatedBy   string `json:"created-by,omitempty"`
}

// UpdateTodoCommand represents a command to update an existing Todo
type UpdateTodoCommand struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    string `json:"priority,omitempty"`
	CategoryID  string `json:"category-id,omitempty"`
}

// CompleteTodoCommand represents a command to mark a Todo as completed
type CompleteTodoCommand struct {
	ID string `json:"id"`
}

// ArchiveTodoCommand represents a command to archive a Todo
type ArchiveTodoCommand struct {
	ID string `json:"id"`
}

// CreateUserCommand represents a command to create a new User
type CreateUserCommand struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

// UpdateUserProfileCommand represents a command to update user profile
type UpdateUserProfileCommand struct {
	ID        string `json:"id"`
	FirstName string `json:"first-name,omitempty"`
	LastName  string `json:"last-name,omitempty"`
	Email     string `json:"email,omitempty"`
}

// PromoteUserCommand represents a command to promote a user to admin
type PromoteUserCommand struct {
	ID string `json:"id"`
}

// SuspendUserCommand represents a command to suspend a user account
type SuspendUserCommand struct {
	ID string `json:"id"`
}

// CreateCategoryCommand represents a command to create a new Category
type CreateCategoryCommand struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color"`
	CreatedBy   string `json:"created-by,omitempty"`
}

// UpdateCategoryCommand represents a command to update a Category
type UpdateCategoryCommand struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}
