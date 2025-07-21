package roleroute

import (
	"reapp/internal/middleware/authmdw"
	"reapp/internal/modules/user/rolehandler"
	"reapp/internal/modules/user/rolerepository"
	"reapp/internal/modules/user/roleservice"
	"reapp/pkg/register"

	"github.com/gin-gonic/gin"
)

type RoleRoute struct{}

func UseRoleRoute() {
	register.ProvideRoute(&RoleRoute{})
}

func (RoleRoute) RegisterRoute(rg *gin.RouterGroup) {
	auth := authmdw.AuthRequired()
	r := rg.Group("/roles").Use(auth)
	h := rolehandler.NewRoleHandler(roleservice.NewRoleService(rolerepository.NewRoleRepository()))

	r.GET("", authmdw.Can("roles.read"), h.List)
	r.GET("/all", authmdw.Can("roles.read"), h.GetAll)
	r.POST("", authmdw.Can("roles.create"), h.Create)
	r.PUT("/:id", authmdw.Can("roles.update"), h.Update)
	r.DELETE("/:id", authmdw.Can("roles.delete"), h.Delete)
	r.GET("/:id", authmdw.Can("roles.detail"), h.GetDetail)
}
