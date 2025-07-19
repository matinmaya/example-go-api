package userservice

import (
	"fmt"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/filterscopes"
	"reapp/pkg/hashcrypto"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type IUserService interface {
	Create(db *gorm.DB, user *usermodel.User) error
	Update(db *gorm.DB, user *usermodel.User) error
	GetByID(db *gorm.DB, id uint32) (*usermodel.User, error)
	Delete(db *gorm.DB, id uint32) error
	List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error
	ChangePassword(db *gorm.DB, data usermodel.ChangePassword) error
}

type UserService struct {
	repository *userrepository.UserRepository
}

func NewUserService(r *userrepository.UserRepository) IUserService {
	return &UserService{repository: r}
}

func (s *UserService) Create(db *gorm.DB, user *usermodel.User) error {
	password, err := hashcrypto.HashMake(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = password
	return s.repository.Create(db, user)
}

func (s *UserService) Update(db *gorm.DB, user *usermodel.User) error {
	if _, err := s.repository.GetByID(db, user.ID); err != nil {
		return fmt.Errorf("something went wrong")
	}

	return s.repository.Update(db, user)
}

func (s *UserService) GetByID(db *gorm.DB, id uint32) (*usermodel.User, error) {
	return s.repository.GetByID(db, id)
}

func (s *UserService) Delete(db *gorm.DB, id uint32) error {
	return s.repository.Delete(db, id)
}

func (s *UserService) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	return s.repository.List(db, pg, filters)
}

func (s *UserService) ChangePassword(db *gorm.DB, data usermodel.ChangePassword) error {
	user, err := s.repository.GetByID(db, data.UserID)
	if err != nil {
		return fmt.Errorf("something went wrong")
	}

	password, err := hashcrypto.HashMake(data.NewPassword)
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	user.Password = password

	return s.repository.Update(db, user)
}
