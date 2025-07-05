package query

// ListTodosQuery represents a query to retrieve all todos following CQRS pattern
type ListTodosQuery struct {
	// Future: Add filtering, pagination, sorting options
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}
