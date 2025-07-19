package register

import (
	"github.com/gin-gonic/gin"
)

type IRouterRegister interface {
	RegisterRoute(rg *gin.RouterGroup)
}

var routes []IRouterRegister

func ProvideRoute(r IRouterRegister) {
	routes = append(routes, r)
}

func InjectRoutes(rg *gin.RouterGroup) {
	for _, rr := range routes {
		rr.RegisterRoute(rg)
	}
}
