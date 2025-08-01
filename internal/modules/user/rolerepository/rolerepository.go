package rolerepository

import (
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) Create(db *gorm.DB, role *rolemodel.Role) error {
	return db.Create(role).Error
}

func (r *RoleRepository) Update(db *gorm.DB, role *rolemodel.Role) error {
	return db.Save(role).Error
}

func (r *RoleRepository) Delete(db *gorm.DB, id uint16) error {
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
	err := db.Order("created_at DESC").Find(&roles).Error

	return roles, err
}

func (r *RoleRepository) List(db *gorm.DB, pg *paginator.Pagination, filters []filterscopes.QueryFilter) error {
	var roles []rolemodel.Role
	scope := paginator.Paginate(db, &rolemodel.Role{}, pg, filters)

	err := db.Scopes(scope).Find(&roles).Error
	if err != nil {
		return err
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
