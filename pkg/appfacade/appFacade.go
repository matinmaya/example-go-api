package appfacade

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/http/reqctx"
	"reapp/pkg/lang"
	"reapp/pkg/services/jwtservice"
)

type ContextFacade struct {
	ctx *gin.Context
}

type AppFacade struct {
	Context    *gin.Context
	DB         *gorm.DB
	Auth       *jwtservice.Claims
	FieldNames []string
}

func NewContextFacade(ctx *gin.Context) *ContextFacade {
	return &ContextFacade{ctx}
}

func (f *ContextFacade) Build() *AppFacade {
	auth, _ := f.ctx.Get("jwt_token")
	db := f.ctx.MustGet("db").(*gorm.DB)
	fls, _ := reqctx.GetFieldNames(f.ctx)

	return &AppFacade{
		Context:    f.ctx,
		DB:         db,
		Auth:       auth.(*jwtservice.Claims),
		FieldNames: fls,
	}
}

func (app *AppFacade) Tran(group string, key string) string {
	return lang.Tran(app.Context, group, key)
}
