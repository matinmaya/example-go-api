package roleroute

import (
	"reapp/internal/middleware/authmw"
	"reapp/internal/modules/user/rolehandler"
	"reapp/internal/modules/user/rolerepository"
	"reapp/internal/modules/user/roleservice"
	"reapp/pkg/http/register"

	"github.com/gin-gonic/gin"
)

type RoleRoute struct{}

func UseRoleRoute() {
	register.AddRoute(&RoleRoute{})
}

func (RoleRoute) RegisterRoute(rg *gin.RouterGroup) {
	auth := authmw.AuthRequired()
	r := rg.Group("/roles").Use(auth)
	h := rolehandler.NewRoleHandler(roleservice.NewRoleService(rolerepository.NewRoleRepository()))

	r.GET("", authmw.Can("roles.read"), h.List)
	r.GET("/all", authmw.Can("roles.read"), h.GetAll)
	r.POST("", authmw.Can("roles.create"), h.Create)
	r.PUT("/:id", authmw.Can("roles.update"), h.Update)
	r.DELETE("/:id", authmw.Can("roles.delete"), h.Delete)
	r.GET("/:id", authmw.Can("roles.detail"), h.GetDetail)
}
