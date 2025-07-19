package userroute

import (
	"reapp/internal/middleware/authmiddleware"
	"reapp/internal/modules/user/userhandler"
	"reapp/internal/modules/user/userrepository"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/register"

	"github.com/gin-gonic/gin"
)

type UserRoute struct{}

func UseUserRoute() {
	register.ProvideRoute(&UserRoute{})
}

func (UserRoute) RegisterRoute(rg *gin.RouterGroup) {
	auth := authmiddleware.AuthRequired()
	r := rg.Group("/users").Use(auth)
	h := userhandler.NewUserHandler(userservice.NewUserService(userrepository.NewUserRepository()))

	r.GET("", authmiddleware.Can("users.read"), h.GetAll)
	r.POST("", authmiddleware.Can("users.create"), h.Create)
	r.PUT("/:id", authmiddleware.Can("users.update"), h.Update)
	r.DELETE("/:id", authmiddleware.Can("users.delete"), h.Delete)
	r.POST("/change-password", authmiddleware.Can("users.change-password"), h.ChangePassword)
}
