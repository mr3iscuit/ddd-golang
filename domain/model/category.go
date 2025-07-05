package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CategoryID represents a unique Category identifier
type CategoryID string

// CategoryColor represents the color theme for a category
type CategoryColor string

const (
	CategoryColorRed    CategoryColor = "red"
	CategoryColorBlue   CategoryColor = "blue"
	CategoryColorGreen  CategoryColor = "green"
	CategoryColorYellow CategoryColor = "yellow"
	CategoryColorPurple CategoryColor = "purple"
	CategoryColorOrange CategoryColor = "orange"
	CategoryColorGray   CategoryColor = "gray"
)

// Category represents a category for organizing todos
type Category struct {
	id          CategoryID
	name        string
	description string
	color       CategoryColor
	createdBy   UserID
	createdAt   time.Time
	updatedAt   time.Time
	isDefault   bool
}

// NewCategory creates a new Category with descriptive factory method
func NewCategory(name string, description string, color CategoryColor, createdBy UserID) *Category {
	now := time.Now()
	return &Category{
		id:          CategoryID(uuid.NewString()),
		name:        name,
		description: description,
		color:       color,
		createdBy:   createdBy,
		createdAt:   now,
		updatedAt:   now,
		isDefault:   false,
	}
}

// NewDefaultCategory creates a default category
func NewDefaultCategory(name string, color CategoryColor) *Category {
	now := time.Now()
	return &Category{
		id:          CategoryID(uuid.NewString()),
		name:        name,
		description: "Default category",
		color:       color,
		createdBy:   "", // System created
		createdAt:   now,
		updatedAt:   now,
		isDefault:   true,
	}
}

// Getters with descriptive names
func (c *Category) GetID() CategoryID {
	return c.id
}

func (c *Category) GetName() string {
	return c.name
}

func (c *Category) GetDescription() string {
	return c.description
}

func (c *Category) GetColor() CategoryColor {
	return c.color
}

func (c *Category) GetCreatedBy() UserID {
	return c.createdBy
}

func (c *Category) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c *Category) GetUpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Category) IsDefault() bool {
	return c.isDefault
}

// Domain behaviors
func (c *Category) UpdateName(newName string) error {
	if newName == "" {
		return errors.New("category name cannot be empty")
	}
	if len(newName) > 50 {
		return errors.New("category name cannot exceed 50 characters")
	}

	c.name = newName
	c.updatedAt = time.Now()
	return nil
}

func (c *Category) UpdateDescription(newDescription string) error {
	if len(newDescription) > 200 {
		return errors.New("category description cannot exceed 200 characters")
	}

	c.description = newDescription
	c.updatedAt = time.Now()
	return nil
}

func (c *Category) UpdateColor(newColor CategoryColor) error {
	switch newColor {
	case CategoryColorRed, CategoryColorBlue, CategoryColorGreen,
		CategoryColorYellow, CategoryColorPurple, CategoryColorOrange, CategoryColorGray:
		c.color = newColor
		c.updatedAt = time.Now()
		return nil
	default:
		return errors.New("invalid category color")
	}
}

func (c *Category) MarkAsDefault() error {
	if c.IsDefault() {
		return errors.New("category is already default")
	}

	c.isDefault = true
	c.updatedAt = time.Now()
	return nil
}

func (c *Category) RemoveDefaultStatus() error {
	if !c.IsDefault() {
		return errors.New("category is not default")
	}

	c.isDefault = false
	c.updatedAt = time.Now()
	return nil
}

// Validation methods
func (c *Category) IsValid() error {
	if c.name == "" {
		return errors.New("category name is required")
	}
	if len(c.name) > 50 {
		return errors.New("category name is too long")
	}
	if len(c.description) > 200 {
		return errors.New("category description is too long")
	}
	return nil
}
