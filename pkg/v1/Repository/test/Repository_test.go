package repository_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yishak-cs/CleanGrpc/Internal/model"
	Repo "github.com/yishak-cs/CleanGrpc/pkg/v1/Repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// if a table exists from a previous run drop it
	db.Migrator().DropTable(&model.User{})

	// migrate the schema
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// create a new user
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
}

func TestRepository_GetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// Create user
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)

	//get user by Id
	id := createdUser.ID
	fetchedUser, err := repo.GetUser(fmt.Sprintf("%d", id))
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Name, fetchedUser.Name)
	assert.Equal(t, createdUser.Email, fetchedUser.Email)

	// Test case: Get non-existent user
	_, err = repo.GetUser("999999")
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepository_GetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// Create a test user first
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Test case: Get user by email
	fetchedUser, err := repo.GetUserByEmail(createdUser.Email)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Name, fetchedUser.Name)
	assert.Equal(t, createdUser.Email, fetchedUser.Email)

	// Test case: Get non-existent email
	_, err = repo.GetUserByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepository_GetUsersList(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// Create multiple test users
	users := []*model.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range users {
		_, err := repo.CreateUser(user)
		assert.NoError(t, err)
	}

	// Test case: Get all users
	usersList := repo.GetUsersList()
	assert.Len(t, usersList, len(users))
}

func TestRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// Create a test user first
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Test case: Update user
	updatedUser := &model.User{
		Model: gorm.Model{ID: createdUser.ID},
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	err = repo.UpdateUser(updatedUser)
	assert.NoError(t, err)

	// Verify the update
	fetchedUser, err := repo.GetUser(fmt.Sprint(createdUser.ID))
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Name, fetchedUser.Name)
	assert.Equal(t, updatedUser.Email, fetchedUser.Email)
}

func TestRepository_DeleteUser(t *testing.T) {
	db := setupTestDB(t)
	repo := Repo.NewRepo(db)

	// Create a test user first
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	createdUser, err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Test case: Delete user
	err = repo.DeleteUser(fmt.Sprint(createdUser.ID))
	assert.NoError(t, err)

	// Verify the deletion
	_, err = repo.GetUser(fmt.Sprint(createdUser.ID))
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
