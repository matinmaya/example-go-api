package provider

import (
	"github.com/go-playground/validator/v10"

	"reapp/internal/modules/customer"
	"reapp/internal/modules/user/usermigration"
	"reapp/pkg/base/basemodel"
	"reapp/pkg/filesystem"
	"reapp/pkg/services/jwtservice"
	"reapp/pkg/validators"
)

func (p *Provider) RegisterServiceProvider() *Provider {
	jwtservice.InitJWT(
		p.cf.JWT.Secret,
		p.cf.JWT.AccessTokenTTL,
		p.cf.JWT.RefreshTokenTTL,
	)

	filesystem.SetPrefixRoutePath(p.cf.Storage.PrefixRoute)
	p.db.AutoMigrate(&basemodel.TableLog{}, &basemodel.HttpLog{})
	usermigration.Migrate(p.db)
	p.db.AutoMigrate(&customer.Customer{})

	customValidator(p)

	return p
}

func customValidator(p *Provider) {
	vlt := validators.InitValidation(p.db, validator.New())
	vlt.RegisterValidation("unique", validators.Unique)
	vlt.RegisterValidation("path", validators.Path)
	vlt.RegisterValidation("slug_strict", validators.SlugStrict)
}
