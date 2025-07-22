package authmdw

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/helpers/jwthelper"
	"reapp/internal/helpers/redishelper"
	"reapp/internal/lang"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/authctx"
	"reapp/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ctxhelper.GetDB(ctx)
		if !jwthelper.ExistsSecret() {
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

		claims, err := jwthelper.ParseToken(parts[1])
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "invalid_token"), nil)
			ctx.Abort()
			return
		}

		redisClient := redishelper.Client()
		if redisClient != nil {
			revoked, _ := redisClient.Get("revoked:" + claims.Id).Result()
			if revoked == "true" {
				response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "token_revoked"), nil)
				ctx.Abort()
				return
			}
		}

		var user usermodel.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "user_not_found"), nil)
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(authctx.SetUserID(ctx.Request.Context(), claims.UserID))
		ctx.Set(ctxhelper.GetDBContextKey(), db.WithContext(ctx.Request.Context()))
		ctx.Set("jwt_token", claims)
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role_ids", claims.RoleIDs)
		ctx.Next()
	}

}
