package repository

import (
	"fmt"

	"github.com/yishak-cs/CleanGrpc/Internal/model"
	"gorm.io/gorm"
)

type Repo struct {
	db gorm.DB
}

func NewRepo(db gorm.DB) *Repo {
	return &Repo{db}
}

func (repo *Repo) CreateUser(user *model.User) (*model.User, error) {
	err := repo.db.Create(user).Error
	if err != nil {
		return &model.User{}, fmt.Errorf("unable to create user: %w", err)
	}
	return user, nil
}

func (repo *Repo) GetUser(id string) (*model.User, error) {
	var user model.User
	if resp := repo.db.First(&user, id).Error; resp != nil {
		return &user, fmt.Errorf("failed to get user: %w", resp)
	}
	return &user, nil
}

func (repo *Repo) GetUsersList() []*model.User {
	var users []*model.User
	resp := repo.db.Find(&users)
	fmt.Printf("%d rows affected", resp.RowsAffected)
	return users
}

func (repo *Repo) UpdateUser(data *model.User) error {
	user, err := repo.GetUser(string(data.ID))
	if err != nil {
		return err
	}
	user.Name = data.Name
	user.Email = data.Email

	if err := repo.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (repo *Repo) DeleteUser(id string) error {
	if err := repo.db.Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete the user: %w", err)
	}

	return nil
}

func (repo *Repo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repo.db.Where("email=?", email).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}
