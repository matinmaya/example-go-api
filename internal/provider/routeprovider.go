package provider

import (
	"reapp/internal/middleware/dbmw"
	"reapp/internal/middleware/langmw"
	"reapp/internal/middleware/logmw"
	"reapp/internal/router"
	"reapp/pkg/http/register"
)

func (p *Provider) RegisterRouteProvider() *Provider {
	p.r.Use(logmw.HttpLogger(p.db))
	p.r.Use(langmw.Language())
	p.r.Use(dbmw.WithDBContext(p.db))

	router.NewRouter(p.r, p.db).UseAdminRouter().UseFrontendRouter().UseNotFoundRouter()
	register.UseRoutes(p.r.Group("/"))
	return p
}
