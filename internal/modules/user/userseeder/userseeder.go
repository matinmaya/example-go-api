package userseeder

import (
	"reapp/internal/modules/user/permmodel"
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/crypto"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	permissions := []permmodel.Permission{
		{Name: "roles.read", Description: "Can read roles"},
		{Name: "roles.detail", Description: "Can read role's detail"},
		{Name: "roles.create", Description: "Can create roles"},
		{Name: "roles.update", Description: "Can update roles"},
		{Name: "roles.delete", Description: "Can delete roles"},
		{Name: "users.read", Description: "Can read users"},
		{Name: "users.create", Description: "Can create users"},
		{Name: "users.update", Description: "Can update users"},
		{Name: "users.delete", Description: "Can delete users"},
		{Name: "permissions.read", Description: "Can read permissions"},
		{Name: "permissions.create", Description: "Can create permissions"},
		{Name: "permissions.update", Description: "Can update permissions"},
		{Name: "permissions.delete", Description: "Can delete permissions"},
	}

	for _, p := range permissions {
		db.FirstOrCreate(&p, permmodel.Permission{Name: p.Name})
	}

	superAdminRole := rolemodel.Role{
		Name:        "superadmin",
		Description: "Super Administrator with full permissions",
		Status:      true,
	}
	db.FirstOrCreate(&superAdminRole, rolemodel.Role{Name: superAdminRole.Name})

	var allPermissions []permmodel.Permission
	db.Find(&allPermissions)

	for _, p := range allPermissions {
		rp := rolemodel.RolePermission{
			RoleID:       superAdminRole.ID,
			PermissionID: p.ID,
		}
		db.FirstOrCreate(&rp, "role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID)
	}

	psw, err := crypto.Make("admin123")
	if err != nil {
		return nil
	}

	admin := usermodel.User{
		Username: "superadmin",
		Password: psw,
		Status:   true,
	}
	db.FirstOrCreate(&admin, usermodel.User{Username: admin.Username})

	ur := usermodel.UserRole{
		UserID: admin.ID,
		RoleID: superAdminRole.ID,
	}
	db.FirstOrCreate(&ur, "user_id = ? AND role_id = ?", ur.UserID, ur.RoleID)

	return nil
}
