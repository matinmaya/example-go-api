package authroute

import (
	"reapp/internal/middleware/authmw"
	"reapp/internal/modules/user/authhandler"
	"reapp/internal/modules/user/authservice"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/http/register"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct{}

func UseAuthRoute() {
	register.AddRoute(&AuthRoute{})
}

func (AuthRoute) RegisterRoute(rg *gin.RouterGroup) {
	r := rg.Group("auth")
	h := authhandler.NewAuthHandler(authservice.NewAuthService(userrepository.NewUserRepository()))

	r.POST("/login", h.Login)
	r.POST("/logout", authmw.AuthRequired(), h.Logout)
	r.POST("/refresh", authmw.RefreshToken(), h.Refresh)
}
