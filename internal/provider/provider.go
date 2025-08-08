package provider

import (
	"reapp/config"
	"reapp/internal/middleware/dbmw"
	"reapp/internal/middleware/langmw"
	"reapp/internal/middleware/logmw"
	"reapp/internal/modules/user/usermigration"
	"reapp/internal/router"
	"reapp/pkg/base/basemodel"
	"reapp/pkg/http/register"
	"reapp/pkg/services/jwtservice"
	"reapp/pkg/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Provider struct {
	r  *gin.Engine
	db *gorm.DB
	cf *config.Config
}

func NewProvider(r *gin.Engine, db *gorm.DB, cf *config.Config) *Provider {
	return &Provider{
		r:  r,
		db: db,
		cf: cf,
	}
}

func (p *Provider) RegisterServiceProvider() *Provider {
	jwtservice.InitJWT(
		p.cf.JWT.Secret,
		p.cf.JWT.AccessTokenTTL,
		p.cf.JWT.RefreshTokenTTL,
	)

	p.db.AutoMigrate(&basemodel.TableLog{}, &basemodel.HttpLog{})
	usermigration.Migrate(p.db)

	vlt := validators.InitValidation(p.db, validator.New())
	vlt.RegisterValidation("unique", validators.Unique)
	vlt.RegisterValidation("path", validators.Path)
	return p
}

func (p *Provider) RegisterRouteProvider() *Provider {
	p.r.Use(logmw.HttpLogger(p.db))
	p.r.Use(langmw.Language())
	p.r.Use(dbmw.WithDBContext(p.db))

	router.NewRouter(p.r, p.db).UseAdminRouter().UseFrontendRouter().UseNotFoundRouter()
	register.UseRoutes(p.r.Group("/"))
	return p
}
