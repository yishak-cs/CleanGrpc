package usecase

import (
	"errors"
	"fmt"

	"github.com/yishak-cs/CleanGrpc/Internal/model"
	interfaces "github.com/yishak-cs/CleanGrpc/pkg/v1"
	"gorm.io/gorm"
)

type UseCase struct {
	repo interfaces.RepoInterface
}

func NewUseCase(repo interfaces.RepoInterface) UseCase {
	return UseCase{repo}
}

func (uc *UseCase) CreateUser(user *model.User) (*model.User, error) {
	if _, err := uc.repo.GetUserByEmail(user.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.User{}, gorm.ErrDuplicatedKey
	}
	return uc.repo.CreateUser(user)
}

func (uc *UseCase) GetUser(id string) (*model.User, error) {
	return uc.repo.GetUser(id)
}

func (uc *UseCase) GetUsersList() []*model.User {
	return uc.repo.GetUsersList()
}

func (uc *UseCase) UpdateUser(update *model.User) error {

	//check if the user exists
	if _, err := uc.repo.GetUser(string((*update).ID)); err != nil {
		return err
	}

	//check if the email is available
	if _, err := uc.repo.GetUserByEmail(update.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("the email already exists. please choose another email")
	}

	// update the user
	if err := uc.repo.UpdateUser(update); err != nil {
		return fmt.Errorf("something went wrong: %w", err)
	}

	return nil
}
func (uc *UseCase) DeleteUser(id string) error {
	var err error
	// check if user exists
	if _, err = uc.GetUser(id); err != nil {
		return err
	}

	err = uc.repo.DeleteUser(id)
	if err != nil {
		// handle the error as it might be something worth to debug
		return err
	}

	return nil
}
