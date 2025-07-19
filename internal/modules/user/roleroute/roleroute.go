package roleroute

import (
	"reapp/internal/middleware/authmiddleware"
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
	auth := authmiddleware.AuthRequired()
	r := rg.Group("/roles").Use(auth)
	h := rolehandler.NewRoleHandler(roleservice.NewRoleService(rolerepository.NewRoleRepository()))

	r.GET("", authmiddleware.Can("roles.read"), h.List)
	r.GET("/all", authmiddleware.Can("roles.read"), h.GetAll)
	r.POST("", authmiddleware.Can("roles.create"), h.Create)
	r.PUT("/:id", authmiddleware.Can("roles.update"), h.Update)
	r.DELETE("/:id", authmiddleware.Can("roles.delete"), h.Delete)
	r.GET("/:id", authmiddleware.Can("roles.detail"), h.GetDetail)
}
