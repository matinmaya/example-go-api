package provider

import (
	"context"
	"log"

	"reapp/config"
	"reapp/internal/file"
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

	storage := fileStorage(p.cf)
	file.RegisterFileRoute(p.r.Group("storages"), storage)

	router.NewRouter(p.r, p.db).UseAdminRouter().UseFrontendRouter().UseNotFoundRouter()
	register.UseRoutes(p.r.Group("/"))
	return p
}

func fileStorage(cf *config.Config) file.Storage {
	var storage file.Storage
	var err error

	switch cf.Storage.Provider {
	case "s3":
		storage, err = file.NewS3StorageFromEnv(context.Background())
		if err != nil {
			log.Fatalf("failed init s3: %v", err)
		}
	default:
		storage, err = file.NewLocalStorage(cf.Storage.Local.BasePath)
		if err != nil {
			log.Fatalf("failed init local storage: %v", err)
		}
	}

	return storage
}
