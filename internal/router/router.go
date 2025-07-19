package router

import (
	"net/http"
	"reapp/internal/modules/user/authroute"
	"reapp/internal/modules/user/roleroute"
	"reapp/internal/modules/user/userroute"
	"reapp/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	return r
}

func (r *Router) UseFrontendRouter() *Router {
	return r
}

func (r *Router) UseNotFoundRouter() *Router {
	r.route.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, http.StatusNotFound, "not found", nil)
	})
	return r
}
