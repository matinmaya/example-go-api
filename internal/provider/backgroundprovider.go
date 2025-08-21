package provider

import (
	"time"

	"reapp/internal/cleaner"
	"reapp/pkg/filesystem"
)

func (p *Provider) RegisterBackgroundProvider() *Provider {
	if filesystem.IsValidStoragePath(p.cf.Storage.Cache.Path) && p.cf.Storage.Cache.CleanupIntervalMin > 0 {
		cleaner := cleaner.NewCleaner(p.cf.Storage.Cache.Path, time.Duration(p.cf.Storage.Cache.MaxAgeMin)*time.Minute)
		cleaner.Start(time.Duration(p.cf.Storage.Cache.CleanupIntervalMin) * time.Minute)
	}

	return p
}
