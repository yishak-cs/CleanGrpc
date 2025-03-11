package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yishak-cs/CleanGrpc/Internal/model"
	usecase "github.com/yishak-cs/CleanGrpc/pkg/v1/UseCase"
	"gorm.io/gorm"
)

// MockRepository is a mock implementation of the RepoInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUser(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUsersList() []*model.User {
	args := m.Called()
	return args.Get(0).([]*model.User)
}

func (m *MockRepository) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUseCase_CreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := usecase.NewUseCase(mockRepo)

	// Test case: Create a new user successfully
	user := &model.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	expectedUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Setup mock expectations
	mockRepo.On("GetUserByEmail", user.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", user).Return(expectedUser, nil)

	// Call the method
	createdUser, err := useCase.CreateUser(user)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, createdUser)
	mockRepo.AssertExpectations(t)

	// Test case: Email already exists
	mockRepo.ExpectedCalls = nil
	existingUser := &model.User{
		Model: gorm.Model{ID: 2},
		Name:  "Existing User",
		Email: "existing@example.com",
	}

	mockRepo.On("GetUserByEmail", existingUser.Email).Return(existingUser, nil)

	// Call the method
	_, err = useCase.CreateUser(existingUser)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := usecase.NewUseCase(mockRepo)

	// Test case: Get existing user
	expectedUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.On("GetUser", "1").Return(expectedUser, nil)

	// Call the method
	user, err := useCase.GetUser("1")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)

	// Test case: User not found
	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetUser", "999").Return(nil, gorm.ErrRecordNotFound)

	// Call the method
	_, err = useCase.GetUser("999")

	// Assertions
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetUsersList(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := usecase.NewUseCase(mockRepo)

	// Test case: Get all users
	expectedUsers := []*model.User{
		{Model: gorm.Model{ID: 1}, Name: "User 1", Email: "user1@example.com"},
		{Model: gorm.Model{ID: 2}, Name: "User 2", Email: "user2@example.com"},
	}

	mockRepo.On("GetUsersList").Return(expectedUsers)

	// Call the method
	users := useCase.GetUsersList()

	// Assertions
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_UpdateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := usecase.NewUseCase(mockRepo)

	// Test case: Update user successfully
	userToUpdate := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	existingUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Original Name",
		Email: "original@example.com",
	}

	mockRepo.On("GetUser", "1").Return(existingUser, nil)
	mockRepo.On("GetUserByEmail", userToUpdate.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("UpdateUser", userToUpdate).Return(nil)

	// Call the method
	err := useCase.UpdateUser(userToUpdate)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test case: User not found
	mockRepo.ExpectedCalls = nil
	nonExistentUser := &model.User{
		Model: gorm.Model{ID: 999},
		Name:  "Non-existent",
		Email: "nonexistent@example.com",
	}

	mockRepo.On("GetUser", "999").Return(nil, gorm.ErrRecordNotFound)

	// Call the method
	err = useCase.UpdateUser(nonExistentUser)

	// Assertions
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	mockRepo.AssertExpectations(t)

	// Test case: Email already exists
	mockRepo.ExpectedCalls = nil
	conflictUser := &model.User{
		Model: gorm.Model{ID: 2},
		Name:  "Conflict User",
		Email: "conflict@example.com",
	}

	existingUser2 := &model.User{
		Model: gorm.Model{ID: 2},
		Name:  "Original Name 2",
		Email: "original2@example.com",
	}

	anotherUser := &model.User{
		Model: gorm.Model{ID: 3},
		Name:  "Another User",
		Email: "conflict@example.com", // Same email as conflictUser
	}

	mockRepo.On("GetUser", "2").Return(existingUser2, nil)
	mockRepo.On("GetUserByEmail", conflictUser.Email).Return(anotherUser, nil)

	// Call the method
	err = useCase.UpdateUser(conflictUser)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUseCase_DeleteUser(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := usecase.NewUseCase(mockRepo)

	// Test case: Delete user successfully
	existingUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.On("GetUser", "1").Return(existingUser, nil)
	mockRepo.On("DeleteUser", "1").Return(nil)

	// Call the method
	err := useCase.DeleteUser("1")

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test case: User not found
	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetUser", "999").Return(nil, gorm.ErrRecordNotFound)

	// Call the method
	err = useCase.DeleteUser("999")

	// Assertions
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	mockRepo.AssertExpectations(t)

	// Test case: Delete operation fails
	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetUser", "2").Return(existingUser, nil)
	mockRepo.On("DeleteUser", "2").Return(errors.New("database error"))

	// Call the method
	err = useCase.DeleteUser("2")

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}
