package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/internal/modules/customer"
	"reapp/internal/modules/user/authroute"
	"reapp/internal/modules/user/roleroute"
	"reapp/internal/modules/user/userroute"
	"reapp/pkg/http/response"
	"reapp/pkg/lang"
)

type Router struct {
	route *gin.Engine
	db    *gorm.DB
}

func NewRouter(r *gin.Engine, db *gorm.DB) *Router {
	return &Router{route: r, db: db}
}

func (r *Router) UseAdminRouter() *Router {
	roleroute.UseRoleRoute()
	userroute.UseUserRoute()
	authroute.UseAuthRoute()
	customer.UseCustomerRoute()
	return r
}

func (r *Router) UseFrontendRouter() *Router {
	return r
}

func (r *Router) UseNotFoundRouter() *Router {
	r.route.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, http.StatusNotFound, lang.Tran(ctx, "response", "not_found"), nil)
	})
	return r
}
