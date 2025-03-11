package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/yishak-cs/CleanGrpc/Internal/model"
	interfaces "github.com/yishak-cs/CleanGrpc/pkg/v1"
	pb "github.com/yishak-cs/CleanGrpc/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type UserServiceServer struct {
	usecase interfaces.UseCaseInterface
	pb.UnimplementedUserServiceServer
}

func NewServer(ser *grpc.Server, uc interfaces.UseCaseInterface) {
	//create an instance of UserServiceServer
	server := UserServiceServer{usecase: uc}
	//register a server which provides the stubs for clients
	pb.RegisterUserServiceServer(ser, &server)
}

func (server *UserServiceServer) CreateUser(ctx context.Context, creq *pb.CreateUserRequest) (*pb.Response, error) {
	//transform the CreatUserRequest type to model.User type
	model := server.transformMessageToModel(creq)

	//validate that the models Name and Email are not empty strings
	if model.Email == "" || model.Name == "" {
		return &pb.Response{Status: "bad request"}, errors.New("please provide your name and email")
	}

	//call UseCase's CreateUser method which accepts User model
	_, err := server.usecase.CreateUser(model)
	if err != nil {
		return &pb.Response{Status: "Something went wrong"}, err
	}

	return &pb.Response{Status: "User Created Successfully"}, nil
}

func (server *UserServiceServer) GetUsersList(ctx context.Context, empty *pb.Empty) (*pb.UsersList, error) {
	//get all the user model instances
	UserList := server.usecase.GetUsersList()

	//create a slice of pointers to UserResponse
	userResponses := []*pb.UserResponse{}

	// loop through the user transforming to and appending UserResponse
	for _, k := range UserList {
		userResponses = append(userResponses, server.transformModelToMessage(k))
	}

	// create the user list response according to the message definition
	users := &pb.UsersList{Users: userResponses}

	return users, nil
}

func (server *UserServiceServer) GetUser(ctx context.Context, req *pb.SingleUserRequest) (*pb.UserResponse, error) {
	//call usecase's GetUser model which accepts id string and return a model instance
	user, err := server.usecase.GetUser(req.Id)

	//handle error
	if err != nil {
		return &pb.UserResponse{}, err
	}

	//transform the model to UserResponse
	userResponse := server.transformModelToMessage(user)

	return userResponse, nil
}

func (server *UserServiceServer) UpdateUser(ctx context.Context, upreq *pb.UpdateUserRequest) (*pb.Response, error) {
	// Create user model from request
	user := &model.User{
		Model: gorm.Model{ID: uint(upreq.Id)},
		Name:  upreq.Name,
		Email: upreq.Email,
	}

	// Call usecase update method
	err := server.usecase.UpdateUser(user)
	if err != nil {
		return &pb.Response{Status: "Failed to update user"}, err
	}

	return &pb.Response{Status: "User updated successfully"}, nil
}

func (server *UserServiceServer) DeleteUser(ctx context.Context, req *pb.SingleUserRequest) (*pb.Response, error) {
	err := server.usecase.DeleteUser(req.Id)
	if err != nil {
		return &pb.Response{Status: "Failed to delete user"}, err
	}

	return &pb.Response{Status: "User deleted successfully"}, nil
}

func (server *UserServiceServer) transformMessageToModel(message *pb.CreateUserRequest) *model.User {
	model := model.User{
		Name:  message.Name,
		Email: message.Email,
	}
	return &model
}

func (server *UserServiceServer) transformModelToMessage(model *model.User) *pb.UserResponse {
	message := pb.UserResponse{
		Id:    fmt.Sprintf("%d", model.ID),
		Name:  model.Name,
		Email: model.Email,
	}
	return &message
}
