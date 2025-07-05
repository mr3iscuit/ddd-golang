package postgres

import (
	"time"

	"gorm.io/gorm"
)

type TodoRecord struct {
	ID          string `gorm:"primaryKey"`
	Title       string
	Description string
	Priority    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` // optional for soft deletes
}

func (TodoRecord) TableName() string {
	return "todos"
}
