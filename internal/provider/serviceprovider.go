package provider

import (
	"github.com/go-playground/validator/v10"

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
	customValidator(p)

	return p
}

func customValidator(p *Provider) {
	vlt := validators.InitValidation(p.db, validator.New())
	vlt.RegisterValidation("unique", validators.Unique)
	vlt.RegisterValidation("path", validators.Path)
	vlt.RegisterValidation("slug_strict", validators.SlugStrict)
}
