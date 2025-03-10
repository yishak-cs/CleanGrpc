package interfaces

import (
	"github.com/yishak-cs/CleanGrpc/Internal/model"
)

type RepoInterface interface {
	CreateUser(model.User) error

	GetUsersList() []*model.User

	GetUser(id string) (*model.User, error)

	UpdateUser(*model.User) error

	DeleteUser(id string) error
}

type UseCaseInterface interface {
	CreateUser(*model.User) error

	GetUsersList() []*model.User

	GetUser(id string) (*model.User, error)

	UpdateUser(*model.User) error

	DeleteUser(id string) error
}
