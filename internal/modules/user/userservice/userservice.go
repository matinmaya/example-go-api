package userservice

import (
	"fmt"
	"log"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/filterscopes"
	"reapp/pkg/hashcrypto"
	"reapp/pkg/lang"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type IUserService interface {
	Create(db *gorm.DB, user *usermodel.User) error
	Update(db *gorm.DB, user *usermodel.User) error
	GetByID(db *gorm.DB, id uint64) (*usermodel.User, error)
	Delete(db *gorm.DB, id uint64) error
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
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "auth", "failed_has_password"))
	}

	tx := db.Begin()
	roleIds := user.RoleIds
	user.Password = password
	if err := s.repository.Create(tx, user); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := s.repository.AsyncUserRoles(tx, user, roleIds); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s *UserService) Update(db *gorm.DB, user *usermodel.User) error {
	if _, err := s.repository.GetByID(db, user.ID); err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	tx := db.Begin()
	roleIds := user.RoleIds
	if err := s.repository.Update(tx, user); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := s.repository.AsyncUserRoles(tx, user, roleIds); err != nil {
		tx.Rollback()
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(tx, "response", "error"))
	}

	return nil
}

func (s *UserService) GetByID(db *gorm.DB, id uint64) (*usermodel.User, error) {
	return s.repository.GetByID(db, uint32(id))
}

func (s *UserService) Delete(db *gorm.DB, id uint64) error {
	if id < 2 {
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	if _, err := s.repository.GetByID(db, uint32(id)); err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "not_found"))
	}

	return s.repository.Delete(db, uint32(id))
}

func (s *UserService) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	return s.repository.List(db, pg, filters)
}

func (s *UserService) ChangePassword(db *gorm.DB, data usermodel.ChangePassword) error {
	user, err := s.repository.GetByID(db, data.UserID)
	if err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "response", "error"))
	}

	password, err := hashcrypto.HashMake(data.NewPassword)
	if err != nil {
		log.Printf("%s", err.Error())
		return fmt.Errorf("%s", lang.TranByDB(db, "auth", "failed_has_password"))
	}
	user.Password = password

	return s.repository.Update(db, user)
}
