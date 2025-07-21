package userrepository

import (
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (UserRepository) Create(db *gorm.DB, user *usermodel.User) error {
	return db.Create(user).Error
}

func (UserRepository) Update(db *gorm.DB, user *usermodel.User) error {
	return db.Save(user).Error
}

func (UserRepository) Delete(db *gorm.DB, id uint32) error {
	return db.Delete(&usermodel.User{}, id).Error
}

func (UserRepository) GetByID(db *gorm.DB, id uint32) (*usermodel.User, error) {
	var user usermodel.User
	err := db.First(&user, id).Error

	return &user, err
}

func (UserRepository) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	var users []usermodel.User
	scope := paginator.Paginate(db, &usermodel.User{}, pg, filters)

	err := db.Scopes(scope).Find(&users).Error
	if err != nil {
		return err
	}

	pg.SetRows(users)
	return nil
}

func (UserRepository) GetByUsername(db *gorm.DB, username string) (*usermodel.User, error) {
	var user usermodel.User
	err := db.Preload("Roles").Where("username =?", username).First(&user).Error

	return &user, err
}

func (UserRepository) SaveTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error {
	return db.Save(tokenInfo).Error
}

func (UserRepository) GetTokenInfoByUserID(db *gorm.DB, userID uint32) (*usermodel.TokenInfo, error) {
	var tokenInfo usermodel.TokenInfo
	err := db.Where("user_id = ?", userID).First(&tokenInfo).Error

	return &tokenInfo, err
}

func (UserRepository) UpdateTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error {
	return db.Save(tokenInfo).Error
}

func (UserRepository) DeleteTokenInfoByUserID(db *gorm.DB, userID uint32) error {
	return db.Where("user_id = ?", userID).Delete(&usermodel.TokenInfo{}).Error
}
