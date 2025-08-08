package authmw

import (
	"net/http"
	"reapp/pkg/http/response"
	"reapp/pkg/lang"
	"reapp/pkg/redisclient"
	"reapp/pkg/services/jwtservice"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		if err == nil {
			redisClient := redisclient.Client()
			if redisClient != nil {
				revoked, _ := redisClient.Get("revoked:" + claims.Id).Result()
				if revoked == "true" {
					response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "token_revoked"), nil)
					ctx.Abort()
					return
				}
			}
			ctx.Set("jwt_token", claims)
		}

		ctx.Next()
	}

}
