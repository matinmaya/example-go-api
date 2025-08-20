package customer

import (
	"github.com/gin-gonic/gin"

	"reapp/internal/middleware/authmw"
	"reapp/pkg/http/register"
)

type CustomerRoute struct{}

func UseCustomerRoute() {
	register.AddRoute(&CustomerRoute{})
}

func (CustomerRoute) RegisterRoute(rg *gin.RouterGroup) {
	auth := authmw.AuthRequired()
	r := rg.Group("/customers").Use(auth)
	h := InitModule()

	r.GET("", authmw.Can("customers.read"), h.List)
	r.POST("", authmw.Can("customers.create"), h.Create)
	r.PUT("/:id", authmw.Can("customers.update"), h.Update)
	r.DELETE("/:id", authmw.Can("customers.delete"), h.Delete)
}
