package rolerepository

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
	"reapp/pkg/services/rediservice"
)

type RoleRepository struct {
	namespace string
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{namespace: "role"}
}

func (r *RoleRepository) Create(db *gorm.DB, role *rolemodel.Role) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	return db.Create(role).Error
}

func (r *RoleRepository) Update(db *gorm.DB, role *rolemodel.Role) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	return db.Save(role).Error
}

func (r *RoleRepository) Delete(db *gorm.DB, id uint16) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	return db.Delete(&rolemodel.Role{}, id).Error
}

func (r *RoleRepository) GetByID(db *gorm.DB, id uint16) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := db.First(&role, id).Error

	return &role, err
}

func (r *RoleRepository) GetDetail(db *gorm.DB, id uint16) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := db.Preload("Permissions").First(&role, id).Error

	return &role, err
}

func (r *RoleRepository) GetAll(db *gorm.DB) ([]rolemodel.Role, error) {
	var roles []rolemodel.Role

	collectionKey := "all"
	if err := rediservice.CacheOfRepository(r.namespace, collectionKey, "data", &roles); err != nil {
		err := db.Order("created_at DESC").Find(&roles).Error
		if err != nil {
			return nil, err
		}

		rediservice.SetCacheOfRepository(r.namespace, collectionKey, "data", roles)
	}

	return roles, nil
}

func (r *RoleRepository) List(ctx *gin.Context, db *gorm.DB, pg *paginator.Pagination, filterFields []queryfilter.FilterField) error {
	var roles []rolemodel.Role
	scopes := paginator.Paginate(db, r.namespace, &rolemodel.Role{}, pg, filterFields)

	collectionKey := "list"
	if err := rediservice.CacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), &roles); err != nil {
		err := db.Scopes(scopes).Find(&roles).Error
		if err != nil {
			return err
		}

		rediservice.SetCacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), roles)
	}

	pg.SetRows(roles)
	return nil
}

func (r *RoleRepository) RoleUserCount(db *gorm.DB, id uint16) (int, error) {
	var count int64
	err := db.Model(&usermodel.UserRole{}).Where("role_id = ?", id).Count(&count).Error
	return int(count), err
}

func (r *RoleRepository) RemovePermissions(db *gorm.DB, roleID uint16) error {
	return db.Where("role_id = ?", roleID).Delete(&rolemodel.RolePermission{}).Error
}

func (r *RoleRepository) AddPermissions(db *gorm.DB, roleID uint16, permissionIDs []uint32) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	var data []rolemodel.RolePermission
	for _, pid := range permissionIDs {
		data = append(data, rolemodel.RolePermission{
			RoleID:       roleID,
			PermissionID: pid,
		})
	}
	return db.Create(&data).Error
}
