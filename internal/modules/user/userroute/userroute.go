package userroute

import (
	"reapp/internal/middleware/authmw"
	"reapp/internal/modules/user/userhandler"
	"reapp/internal/modules/user/userrepository"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/http/register"

	"github.com/gin-gonic/gin"
)

type UserRoute struct{}

func UseUserRoute() {
	register.AddRoute(&UserRoute{})
}

func (UserRoute) RegisterRoute(rg *gin.RouterGroup) {
	auth := authmw.AuthRequired()
	r := rg.Group("/users").Use(auth)
	h := userhandler.NewUserHandler(userservice.NewUserService(userrepository.NewUserRepository()))

	r.GET("", authmw.Can("users.read"), h.List)
	r.POST("", authmw.Can("users.create"), h.Create)
	r.PUT("/:id", authmw.Can("users.update"), h.Update)
	r.DELETE("/:id", authmw.Can("users.delete"), h.Delete)
	r.POST("/change-password", authmw.Can("users.change-password"), h.ChangePassword)
}
