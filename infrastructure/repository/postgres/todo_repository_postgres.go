package postgres

import (
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/mr3iscuit/ddd-golang/application/port"
	"github.com/mr3iscuit/ddd-golang/domain/model"
)

// PostgresTodoRepository implements port.TodoRepositoryPort using PostgreSQL and GORM
type PostgresTodoRepository struct {
	db *gorm.DB
}

// NewPostgresTodoRepository creates a new PostgresTodoRepository
func NewPostgresTodoRepository(db *gorm.DB) *PostgresTodoRepository {
	return &PostgresTodoRepository{db: db}
}

var _ port.TodoRepositoryPort = (*PostgresTodoRepository)(nil)

// Save inserts or updates a Todo in the database
func (r *PostgresTodoRepository) Save(todo *model.Todo) error {
	record := fromModel(todo)
	result := r.db.Save(record)
	return result.Error
}

// FindByID retrieves a Todo by ID
func (r *PostgresTodoRepository) FindByID(id model.TodoID) (*model.Todo, error) {
	var record TodoRecord
	result := r.db.Where("id = ?", id).First(&record)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("todo with id %s not found", id)
		}
		return nil, result.Error
	}
	return toModel(&record), nil
}

// FindAll retrieves all Todos
func (r *PostgresTodoRepository) FindAll() ([]*model.Todo, error) {
	var records []TodoRecord
	result := r.db.Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	todos := make([]*model.Todo, len(records))
	for i := range records {
		todos[i] = toModel(&records[i])
	}
	return todos, nil
}

// Delete removes a Todo by ID
func (r *PostgresTodoRepository) Delete(id model.TodoID) error {
	result := r.db.Delete(&TodoRecord{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}
	return nil
}
