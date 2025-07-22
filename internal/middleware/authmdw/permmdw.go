package authmdw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/helpers/redishelper"
	"reapp/internal/lang"
	"reapp/internal/modules/user/rolemodel"
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
		cacheKey := "permissions:user:" + userID

		var permissions []string
		redisClient := redishelper.Client()
		cached, err := redisClient.Get(cacheKey).Result()
		if err == nil {
			json.Unmarshal([]byte(cached), &permissions)
		} else {
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

			bytes, _ := json.Marshal(permissions)
			redisClient.Set(cacheKey, bytes, time.Hour)
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
