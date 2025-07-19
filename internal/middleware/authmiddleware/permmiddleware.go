package authmiddleware

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/modules/user/rolemodel"
	"reapp/pkg/response"

	"github.com/gin-gonic/gin"
)

func Can(requiredPms string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ctxhelper.GetDB(ctx)
		roleIDsValue, exists := ctx.Get("role_ids")
		if !exists {
			response.Error(ctx, http.StatusUnauthorized, "Role information not found", nil)
			ctx.Abort()
			return
		}

		roleIDs, ok := roleIDsValue.([]uint16)
		if !ok || len(roleIDs) == 0 {
			response.Error(ctx, http.StatusUnauthorized, "Invalid role data", nil)
			ctx.Abort()
			return
		}

		var roles []rolemodel.Role
		if err := db.Preload("Permissions").Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			response.Error(ctx, http.StatusForbidden, "Roles not found", nil)
			ctx.Abort()
			return
		}

		hasPermission := false
		for _, role := range roles {
			for _, permission := range role.Permissions {
				if permission.Name == requiredPms {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response.Error(ctx, http.StatusForbidden, "Insufficient permissions", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
