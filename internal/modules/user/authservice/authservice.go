package authservice

import (
	"fmt"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/hashcrypto"

	"gorm.io/gorm"
)

type IAuthService interface {
	Attempt(db *gorm.DB, cdt usermodel.AuthCredentials) (*usermodel.User, error)
	// Profile(user *usermodel.User)
}

type AuthService struct {
	repository *userrepository.UserRepository
}

func NewAuthService(r *userrepository.UserRepository) IAuthService {
	return &AuthService{repository: r}
}

func (s *AuthService) Attempt(db *gorm.DB, cdt usermodel.AuthCredentials) (*usermodel.User, error) {
	user, err := s.repository.GetByUsername(db, cdt.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !hashcrypto.HashCheck(cdt.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
