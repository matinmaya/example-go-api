package customer

import (
	"reapp/pkg/module"
)

func InitModule() module.IHandler[Customer, uint32, Customer, CustomerListQuery] {
	hYields := module.HandlerYields[Customer, Customer]{
		UpdateValidateScope: updateValidateScope(),
		BeforeResponse:      beforeResponse(),
		BeforeResponseList:  beforeResponseList(),
	}
	sYields := module.ServiceYields[Customer, uint32]{}
	mod := module.NewModule[Customer, uint32, Customer, CustomerListQuery](
		"customers",
		hYields,
		sYields,
	)

	return mod.Handler
}
