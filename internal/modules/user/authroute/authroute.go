package authroute

import (
	"reapp/internal/modules/user/authhandler"
	"reapp/internal/modules/user/authservice"
	"reapp/internal/modules/user/userrepository"
	"reapp/pkg/register"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct{}

func UseAuthRoute() {
	register.ProvideRoute(&AuthRoute{})
}

func (AuthRoute) RegisterRoute(rg *gin.RouterGroup) {
	r := rg.Group("auth")
	h := authhandler.NewAuthHandler(authservice.NewAuthService(userrepository.NewUserRepository()))

	r.POST("/login", h.Login)
}
