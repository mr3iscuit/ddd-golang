package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserID represents a unique User identifier
type UserID string

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// User represents a user in the domain
type User struct {
	id          UserID
	email       string
	username    string
	firstName   string
	lastName    string
	role        UserRole
	status      UserStatus
	createdAt   time.Time
	updatedAt   time.Time
	lastLoginAt *time.Time
}

// NewUser creates a new User with descriptive factory method
func NewUser(email string, username string, firstName string, lastName string) *User {
	now := time.Now()
	return &User{
		id:          UserID(uuid.NewString()),
		email:       email,
		username:    username,
		firstName:   firstName,
		lastName:    lastName,
		role:        UserRoleUser,
		status:      UserStatusActive,
		createdAt:   now,
		updatedAt:   now,
		lastLoginAt: nil,
	}
}

// NewAdminUser creates a new admin user
func NewAdminUser(email string, username string, firstName string, lastName string) *User {
	user := NewUser(email, username, firstName, lastName)
	user.role = UserRoleAdmin
	return user
}

// Getters with descriptive names
func (u *User) GetID() UserID {
	return u.id
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetFirstName() string {
	return u.firstName
}

func (u *User) GetLastName() string {
	return u.lastName
}

func (u *User) GetFullName() string {
	return u.firstName + " " + u.lastName
}

func (u *User) GetRole() UserRole {
	return u.role
}

func (u *User) GetStatus() UserStatus {
	return u.status
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) GetLastLoginAt() *time.Time {
	return u.lastLoginAt
}

// Status checks
func (u *User) IsActive() bool {
	return u.status == UserStatusActive
}

func (u *User) IsAdmin() bool {
	return u.role == UserRoleAdmin
}

func (u *User) IsSuspended() bool {
	return u.status == UserStatusSuspended
}

// Domain behaviors
func (u *User) UpdateProfile(firstName string, lastName string) error {
	if firstName == "" || lastName == "" {
		return errors.New("first name and last name cannot be empty")
	}

	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdateEmail(newEmail string) error {
	if newEmail == "" {
		return errors.New("email cannot be empty")
	}
	// Add email validation logic here if needed

	u.email = newEmail
	u.updatedAt = time.Now()
	return nil
}

func (u *User) PromoteToAdmin() error {
	if u.IsAdmin() {
		return errors.New("user is already an admin")
	}

	u.role = UserRoleAdmin
	u.updatedAt = time.Now()
	return nil
}

func (u *User) DemoteToUser() error {
	if !u.IsAdmin() {
		return errors.New("user is not an admin")
	}

	u.role = UserRoleUser
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ActivateAccount() error {
	if u.IsActive() {
		return errors.New("account is already active")
	}

	u.status = UserStatusActive
	u.updatedAt = time.Now()
	return nil
}

func (u *User) SuspendAccount() error {
	if u.IsSuspended() {
		return errors.New("account is already suspended")
	}

	u.status = UserStatusSuspended
	u.updatedAt = time.Now()
	return nil
}

func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

func (u *User) GetDaysSinceLastLogin() (int, error) {
	if u.lastLoginAt == nil {
		return 0, errors.New("user has never logged in")
	}

	days := int(time.Since(*u.lastLoginAt).Hours() / 24)
	return days, nil
}
