package provider

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/config"
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
