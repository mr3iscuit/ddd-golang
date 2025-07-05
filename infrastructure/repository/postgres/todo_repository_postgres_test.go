package postgres

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mr3iscuit/ddd-golang/domain/model"
)

type PostgresRepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *PostgresTodoRepository
}

func (s *PostgresRepoTestSuite) SetupSuite() {
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost user=todo_user password=todo_password dbname=todo_db port=5432 sslmode=disable"
	}

	var err error
	s.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)

	// Auto-migrate the schema
	err = s.db.AutoMigrate(&TodoRecord{})
	s.Require().NoError(err)

	s.repo = NewPostgresTodoRepository(s.db)
}

func (s *PostgresRepoTestSuite) TearDownTest() {
	// Clear all rows after each test
	s.db.Exec("DELETE FROM todos")
}

func (s *PostgresRepoTestSuite) TestSaveAndFindByID() {
	todo := model.NewTodo("Test Title", "Test Description", model.TodoPriorityHigh)
	err := s.repo.Save(todo)
	s.NoError(err)

	found, err := s.repo.FindByID(todo.GetID())
	s.NoError(err)
	s.Equal(todo.GetID(), found.GetID())
	s.Equal(todo.GetTitle(), found.GetTitle())
	s.Equal(todo.GetDescription(), found.GetDescription())
	s.Equal(todo.GetPriority(), found.GetPriority())
	s.Equal(todo.GetStatus(), found.GetStatus())
	s.WithinDuration(todo.GetCreatedAt(), found.GetCreatedAt(), time.Second)
	s.WithinDuration(todo.GetUpdatedAt(), found.GetUpdatedAt(), time.Second)
}

func (s *PostgresRepoTestSuite) TestFindAll() {
	t1 := model.NewTodo("First", "Desc1", model.TodoPriorityLow)
	t2 := model.NewTodo("Second", "Desc2", model.TodoPriorityMedium)

	s.NoError(s.repo.Save(t1))
	s.NoError(s.repo.Save(t2))

	all, err := s.repo.FindAll()
	s.NoError(err)
	s.Len(all, 2)

	var ids []model.TodoID
	for _, t := range all {
		ids = append(ids, t.GetID())
	}
	s.Contains(ids, t1.GetID())
	s.Contains(ids, t2.GetID())
}

func (s *PostgresRepoTestSuite) TestDelete() {
	todo := model.NewTodo("To be deleted", "", model.TodoPriorityLow)
	s.NoError(s.repo.Save(todo))

	err := s.repo.Delete(todo.GetID())
	s.NoError(err)

	_, err = s.repo.FindByID(todo.GetID())
	s.Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *PostgresRepoTestSuite) TestMarkAsCompleted() {
	todo := model.NewTodo("Complete Me", "", model.TodoPriorityMedium)
	s.NoError(s.repo.Save(todo))

	s.NoError(todo.MarkAsCompleted())
	s.NoError(s.repo.Save(todo))

	found, err := s.repo.FindByID(todo.GetID())
	s.NoError(err)
	s.Equal(model.TodoStatusCompleted, found.GetStatus())
	s.NotNil(found.GetCompletedAt())
	s.WithinDuration(time.Now(), *found.GetCompletedAt(), time.Second)
}

func (s *PostgresRepoTestSuite) TestArchiveTodo() {
	todo := model.NewTodo("Archive Me", "", model.TodoPriorityHigh)
	s.NoError(s.repo.Save(todo))

	s.NoError(todo.ArchiveTodo())
	s.NoError(s.repo.Save(todo))

	found, err := s.repo.FindByID(todo.GetID())
	s.NoError(err)
	s.Equal(model.TodoStatusArchived, found.GetStatus())
}

func TestPostgresRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresRepoTestSuite))
}
