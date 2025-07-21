package authservice

import (
	"fmt"
	"reapp/internal/helpers/redishelper"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/hashcrypto"
	"time"

	"gorm.io/gorm"
)

type IAuthService interface {
	GetUserByID(db *gorm.DB, id uint32) (*usermodel.User, error)
	Attempt(db *gorm.DB, cdt usermodel.AuthCredentials) (*usermodel.User, error)
	SaveTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error
	GetTokenInfoByUserID(db *gorm.DB, userID uint32) (*usermodel.TokenInfo, error)
	UpdateTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error
	DeleteTokenInfoByUserID(db *gorm.DB, userID uint32) error
	RevokeAccessToken(jti string, expiresAt time.Time) error
}

type AuthService struct {
	repository *userrepository.UserRepository
}

func NewAuthService(r *userrepository.UserRepository) IAuthService {
	return &AuthService{repository: r}
}

func (s *AuthService) GetUserByID(db *gorm.DB, id uint32) (*usermodel.User, error) {
	user, err := s.repository.GetByID(db, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
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

func (s *AuthService) SaveTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error {
	return s.repository.SaveTokenInfo(db, tokenInfo)
}

func (s *AuthService) GetTokenInfoByUserID(db *gorm.DB, userID uint32) (*usermodel.TokenInfo, error) {
	return s.repository.GetTokenInfoByUserID(db, userID)
}

func (s *AuthService) UpdateTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error {
	return s.repository.UpdateTokenInfo(db, tokenInfo)
}

func (s *AuthService) DeleteTokenInfoByUserID(db *gorm.DB, userID uint32) error {
	return s.repository.DeleteTokenInfoByUserID(db, userID)
}

func (s *AuthService) RevokeAccessToken(jti string, expiresAt time.Time) error {
	client := redishelper.Client()
	key := "revoked:" + jti
	return client.Set(key, "true", time.Until(expiresAt)).Err()
}
