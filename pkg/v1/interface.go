package interfaces

import (
	"github.com/yishak-cs/CleanGrpc/Internal/model"
)

type RepoInterface interface {
	CreateUser(*model.User) (*model.User, error)

	GetUsersList() []*model.User

	GetUser(id string) (*model.User, error)

	UpdateUser(*model.User) error

	DeleteUser(string) error

	GetUserByEmail(string) (*model.User, error)
}

type UseCaseInterface interface {
	CreateUser(*model.User) (*model.User, error)

	GetUsersList() []*model.User

	GetUser(id string) (*model.User, error)

	UpdateUser(*model.User) error

	DeleteUser(id string) error
}
