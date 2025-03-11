package handler_test

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yishak-cs/CleanGrpc/Internal/model"
	interfaces "github.com/yishak-cs/CleanGrpc/pkg/v1"
	handler "github.com/yishak-cs/CleanGrpc/pkg/v1/handler/grpc"
	pb "github.com/yishak-cs/CleanGrpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

// MockUseCase is a mock implementation of the UseCaseInterface
type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) CreateUser(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUseCase) GetUser(id string) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUseCase) GetUsersList() []*model.User {
	args := m.Called()
	return args.Get(0).([]*model.User)
}

func (m *MockUseCase) UpdateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUseCase) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Fixed setupGrpcServer function that doesn't call t.Fatalf in a goroutine
func setupGrpcServer(t *testing.T, mockUseCase interfaces.UseCaseInterface) (*grpc.ClientConn, pb.UserServiceClient) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	// Register our service
	handler.NewUserServer(s, mockUseCase)

	// Use a WaitGroup to ensure server is started before proceeding
	var wg sync.WaitGroup
	wg.Add(1)

	// Start server in a goroutine
	go func() {
		wg.Done() // Signal that server is ready to accept connections
		if err := s.Serve(lis); err != nil {
			// Don't call t.Fatalf here - just log the error
			// The test will fail naturally if connections can't be established
		}
	}()

	// Wait for server to start
	wg.Wait()

	// Create a client connection
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// Connect to the server - do error handling in the main test thread
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return conn, client
}

func TestUserServiceServer_CreateUser(t *testing.T) {
	mockUseCase := new(MockUseCase)
	conn, client := setupGrpcServer(t, mockUseCase)
	defer conn.Close()

	// Test case: Create user successfully
	createReq := &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	expectedUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Setup mock expectations
	mockUseCase.On("CreateUser", mock.MatchedBy(func(u *model.User) bool {
		return u.Name == createReq.Name && u.Email == createReq.Email
	})).Return(expectedUser, nil)

	// Call the method
	resp, err := client.CreateUser(context.Background(), createReq)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "User Created Successfully", resp.Status)
	mockUseCase.AssertExpectations(t)

	// Test case: Missing required fields
	mockUseCase.ExpectedCalls = nil
	emptyReq := &pb.CreateUserRequest{
		Name:  "",
		Email: "",
	}

	// Call the method - this should return an error but not panic
	_, err = client.CreateUser(context.Background(), emptyReq)

	// Assertions - we expect an error but the test shouldn't crash
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please provide your name and email")
	mockUseCase.AssertNotCalled(t, "CreateUser")

	// Test case: UseCase returns an error
	mockUseCase.ExpectedCalls = nil
	errorReq := &pb.CreateUserRequest{
		Name:  "Error User",
		Email: "error@example.com",
	}

	mockUseCase.On("CreateUser", mock.MatchedBy(func(u *model.User) bool {
		return u.Name == errorReq.Name && u.Email == errorReq.Email
	})).Return(nil, errors.New("database error"))

	// Call the method
	_, err = client.CreateUser(context.Background(), errorReq)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockUseCase.AssertExpectations(t)
}

func TestUserServiceServer_GetUsersList(t *testing.T) {
	mockUseCase := new(MockUseCase)
	conn, client := setupGrpcServer(t, mockUseCase)
	defer conn.Close()

	// Test case: Get users list successfully
	expectedUsers := []*model.User{
		{Model: gorm.Model{ID: 1}, Name: "User 1", Email: "user1@example.com"},
		{Model: gorm.Model{ID: 2}, Name: "User 2", Email: "user2@example.com"},
	}

	mockUseCase.On("GetUsersList").Return(expectedUsers)

	// Call the method
	resp, err := client.GetUsersList(context.Background(), &pb.Empty{})

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, "1", resp.Users[0].Id)
	assert.Equal(t, "User 1", resp.Users[0].Name)
	assert.Equal(t, "user1@example.com", resp.Users[0].Email)
	assert.Equal(t, "2", resp.Users[1].Id)
	mockUseCase.AssertExpectations(t)
}

func TestUserServiceServer_GetUser(t *testing.T) {
	mockUseCase := new(MockUseCase)
	conn, client := setupGrpcServer(t, mockUseCase)
	defer conn.Close()

	// Test case: Get user successfully
	expectedUser := &model.User{
		Model: gorm.Model{ID: 1},
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockUseCase.On("GetUser", "1").Return(expectedUser, nil)

	// Call the method
	resp, err := client.GetUser(context.Background(), &pb.SingleUserRequest{Id: "1"})

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.Id)
	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "test@example.com", resp.Email)
	mockUseCase.AssertExpectations(t)

	// Test case: User not found
	mockUseCase.ExpectedCalls = nil
	mockUseCase.On("GetUser", "999").Return(nil, gorm.ErrRecordNotFound)

	// Call the method
	_, err = client.GetUser(context.Background(), &pb.SingleUserRequest{Id: "999"})

	// Assertions
	assert.Error(t, err)
	mockUseCase.AssertExpectations(t)
}

func TestUserServiceServer_UpdateUser(t *testing.T) {
	mockUseCase := new(MockUseCase)
	conn, client := setupGrpcServer(t, mockUseCase)
	defer conn.Close()

	// Test case: Update user successfully
	updateReq := &pb.UpdateUserRequest{
		Id:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	mockUseCase.On("UpdateUser", mock.MatchedBy(func(u *model.User) bool {
		return u.ID == uint(updateReq.Id) &&
			u.Name == updateReq.Name &&
			u.Email == updateReq.Email
	})).Return(nil)

	// Call the method
	resp, err := client.UpdateUser(context.Background(), updateReq)

	// Assertions - check for nil before accessing
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil")
	if resp != nil {
		assert.Equal(t, "User updated successfully", resp.Status)
	}
	mockUseCase.AssertExpectations(t)

	// Test case: Update fails
	mockUseCase.ExpectedCalls = nil
	errorReq := &pb.UpdateUserRequest{
		Id:    999,
		Name:  "Error User",
		Email: "error@example.com",
	}

	mockUseCase.On("UpdateUser", mock.MatchedBy(func(u *model.User) bool {
		return u.ID == uint(errorReq.Id) &&
			u.Name == errorReq.Name &&
			u.Email == errorReq.Email
	})).Return(errors.New("user not found"))

	// Call the method
	resp, err = client.UpdateUser(context.Background(), errorReq)

	// Assertions - handle potential nil response
	assert.Error(t, err)
	// Only check the response status if resp is not nil
	if resp != nil {
		assert.Equal(t, "Failed to update user", resp.Status)
	}
	mockUseCase.AssertExpectations(t)
}
