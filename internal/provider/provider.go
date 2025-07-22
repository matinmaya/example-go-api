package provider

import (
	"reapp/config"
	"reapp/internal/helpers/jwthelper"
	"reapp/internal/middleware/dbmdw"
	"reapp/internal/middleware/langmdw"
	"reapp/internal/middleware/loggermdw"
	"reapp/internal/models"
	"reapp/internal/modules/user/usermigration"
	"reapp/internal/router"
	"reapp/pkg/basemodel"
	"reapp/pkg/register"
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
	jwthelper.InitJWT(
		p.cf.JWT.Secret,
		p.cf.JWT.AccessTokenTTL,
		p.cf.JWT.RefreshTokenTTL,
	)

	p.db.AutoMigrate(&basemodel.SysLog{}, &models.RequestLog{})
	usermigration.Migrate(p.db)

	vlt := validators.InitValidation(p.db, validator.New())
	vlt.RegisterValidation("unique", validators.Unique)
	vlt.RegisterValidation("path", validators.Path)
	return p
}

func (p *Provider) RegisterRouteProvider() *Provider {
	p.r.Use(loggermdw.RequestLogger(p.db))
	p.r.Use(langmdw.Language())
	p.r.Use(dbmdw.WithDBContext(p.db))

	router.NewRouter(p.r, p.db).UseAdminRouter().UseFrontendRouter().UseNotFoundRouter()
	register.InjectRoutes(p.r.Group("/"))
	return p
}
