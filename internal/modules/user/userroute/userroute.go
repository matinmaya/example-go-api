package userroute

import (
	"reapp/internal/middleware/authmdw"
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
	auth := authmdw.AuthRequired()
	r := rg.Group("/users").Use(auth)
	h := userhandler.NewUserHandler(userservice.NewUserService(userrepository.NewUserRepository()))

	r.GET("", authmdw.Can("users.read"), h.List)
	r.POST("", authmdw.Can("users.create"), h.Create)
	r.PUT("/:id", authmdw.Can("users.update"), h.Update)
	r.DELETE("/:id", authmdw.Can("users.delete"), h.Delete)
	r.POST("/change-password", authmdw.Can("users.change-password"), h.ChangePassword)
}
