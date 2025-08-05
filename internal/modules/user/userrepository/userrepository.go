package userrepository

import (
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/database/redisdb"
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type UserRepository struct {
	namespace string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{namespace: "user"}
}

func (r *UserRepository) Create(db *gorm.DB, user *usermodel.User) error {
	go redisdb.ClearCacheOfRepository(r.namespace)
	return db.Create(user).Error
}

func (r *UserRepository) Update(db *gorm.DB, user *usermodel.User) error {
	go redisdb.ClearCacheOfRepository(r.namespace)
	return db.Save(user).Error
}

func (r *UserRepository) Delete(db *gorm.DB, id uint32) error {
	go redisdb.ClearCacheOfRepository(r.namespace)
	return db.Delete(&usermodel.User{}, id).Error
}

func (UserRepository) GetByID(db *gorm.DB, id uint32) (*usermodel.User, error) {
	var user usermodel.User
	err := db.First(&user, id).Error

	return &user, err
}

func (r *UserRepository) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	var users []usermodel.User
	scope := paginator.Paginate(db, r.namespace, &usermodel.User{}, pg, filters)

	collectionKey := "list"
	if err := redisdb.GetCacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), &users); err != nil {
		err := db.Scopes(scope).Find(&users).Error
		if err != nil {
			return err
		}

		redisdb.SetCacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), users)
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

	var data []usermodel.UserRole
	for _, roleId := range roleIds {
		data = append(data, usermodel.UserRole{
			UserID: user.ID,
			RoleID: roleId,
		})
	}

	return db.Create(&data).Error
}
