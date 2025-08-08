package authmw

import (
	"fmt"
	"net/http"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/context/authctx"
	"reapp/pkg/context/dbctx"
	"reapp/pkg/http/response"
	"reapp/pkg/lang"
	"reapp/pkg/redisclient"
	"reapp/pkg/services/jwtservice"
	"reapp/pkg/services/rediservice"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := dbctx.DB(ctx)
		if !jwtservice.ExistsSecret() {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "jwt_secret_not_set"), nil)
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "authorization_required"), nil)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "authorization_invalid_format"), nil)
			ctx.Abort()
			return
		}

		claims, err := jwtservice.ParseToken(parts[1])
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "invalid_token"), nil)
			ctx.Abort()
			return
		}

		redisClient := redisclient.Client()
		if redisClient != nil {
			revoked, _ := redisClient.Get("revoked:" + claims.Id).Result()
			if revoked == "true" {
				response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "token_revoked"), nil)
				ctx.Abort()
				return
			}
		}

		authID := fmt.Sprintf("%d", claims.UserID)
		var user usermodel.User
		if err := rediservice.CacheOfAuthUser(authID, &user); err != nil {
			if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
				response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "user_not_found"), nil)
				ctx.Abort()
				return
			}
			if err := rediservice.SetCacheOfAuthUser(authID, user); err != nil {
				response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "cache_error"), nil)
				ctx.Abort()
				return
			}
		}

		if !user.Status {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "account_locked"), nil)
			ctx.Abort()
			return
		}

		ctxValue := authctx.SetUserID(ctx.Request.Context(), claims.UserID)
		ctx.Request = ctx.Request.WithContext(ctxValue)
		ctx.Set(dbctx.DBContextKey(), db.WithContext(ctx.Request.Context()))
		ctx.Set("jwt_token", claims)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role_ids", claims.RoleIDs)
		ctx.Next()
	}

}
