package authmdw

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/helpers/jwthelper"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/authctx"
	"reapp/pkg/response"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := ctxhelper.GetDB(ctx)
		if len(jwthelper.GetSecret()) == 0 {
			response.Error(ctx, http.StatusUnauthorized, "JWT secret is not set", nil)
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "Authorization header is required", nil)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, "Invalid authorization header format", nil)
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &jwthelper.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwthelper.GetSecret(), nil
		})

		if err != nil || !token.Valid {
			response.Error(ctx, http.StatusUnauthorized, "Invalid token", nil)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(*jwthelper.Claims)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "Invalid token claims", nil)
			ctx.Abort()
			return
		}

		var user usermodel.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			response.Error(ctx, http.StatusUnauthorized, "User not found", nil)
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(authctx.SetUserID(ctx.Request.Context(), claims.UserID))
		ctx.Set(ctxhelper.GetCxtDBKey(), db.WithContext(ctx.Request.Context()))
		ctx.Set("user_id", claims.RoleIDs)
		ctx.Set("role_ids", claims.RoleIDs)
		ctx.Next()
	}

}
