package authmdw

import (
	"fmt"
	"net/http"

	"reapp/internal/modules/user/rolemodel"
	"reapp/pkg/database/redisdb"
	"reapp/pkg/helpers/ctxhelper"
	"reapp/pkg/lang"
	"reapp/pkg/response"

	"github.com/gin-gonic/gin"
)

func Can(requiredPms string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ctxhelper.GetDB(ctx)
		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "permission", "user_id_not_found"), nil)
			ctx.Abort()
			return
		}

		userID := fmt.Sprintf("%v", userIDValue)
		permissions, err := redisdb.GetCacheOfPerms(userID)
		if err != nil {
			roleIDsValue, _ := ctx.Get("role_ids")
			roleIDs, ok := roleIDsValue.([]uint16)
			if !ok || len(roleIDs) == 0 {
				response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "permission", "invalid_role_data"), nil)
				ctx.Abort()
				return
			}

			var roles []rolemodel.Role
			if err := db.Preload("Permissions").Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
				response.Error(ctx, http.StatusForbidden, lang.Tran(ctx, "permission", "roles_not_found"), nil)
				ctx.Abort()
				return
			}

			permMap := map[string]struct{}{}
			for _, role := range roles {
				for _, p := range role.Permissions {
					permMap[p.Name] = struct{}{}
				}
			}
			for name := range permMap {
				permissions = append(permissions, name)
			}

			if err := redisdb.SetCacheOfPerms(userID, permissions); err != nil {
				response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "permission", "cache_error"), nil)
				ctx.Abort()
				return
			}
		}

		hasPermission := false
		for _, p := range permissions {
			if p == requiredPms {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Error(ctx, http.StatusForbidden, lang.Tran(ctx, "permission", "insufficient_permission"), nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
