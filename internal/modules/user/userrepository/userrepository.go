package userrepository

import (
	"fmt"
	"reapp/internal/modules/user/rolemodel"
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

func (UserRepository) GetTokenInfo(db *gorm.DB, userID uint32, jti string) (*usermodel.TokenInfo, error) {
	var tokenInfo usermodel.TokenInfo
	err := db.Where("user_id = ?", userID).Where("jti = ?", jti).First(&tokenInfo).Error

	return &tokenInfo, err
}

func (UserRepository) UpdateTokenInfo(db *gorm.DB, tokenInfo *usermodel.TokenInfo) error {
	return db.Save(tokenInfo).Error
}

func (UserRepository) DeleteTokenInfo(db *gorm.DB, userID uint32, jti string) error {
	return db.Where("user_id = ?", userID).Where("jti = ?", jti).Delete(&usermodel.TokenInfo{}).Error
}

func (UserRepository) AsyncUserRoles(db *gorm.DB, user *usermodel.User, roleIds []uint16) error {
	if err := db.Model(user).Association("Roles").Clear(); err != nil {
		return err
	}

	if len(roleIds) == 0 {
		return nil
	}

	var roles []rolemodel.Role
	if err := db.Where("id IN ?", roleIds).Find(&roles).Error; err != nil {
		return err
	}

	fmt.Printf("roles count %v\n", db.Model(user).Association("Roles").Count())

	// error here
	if err := db.Model(user).Association("Roles").Replace(&roles); err != nil {
		return err
	}

	return nil
}
