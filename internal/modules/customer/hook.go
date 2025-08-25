package customer

import (
	"github.com/gin-gonic/gin"

	"reapp/pkg/base/basehandler"
	"reapp/pkg/filesystem"
	"reapp/pkg/validators"
)

func updateValidateScope() basehandler.TWithIDScope[Customer] {
	return func(ctx *gin.Context, cus *Customer, id uint64) error {
		cus.UniqueScope = validators.ExceptByID(id)
		return nil
	}
}

func beforeResponseList() basehandler.TResponseListHook[Customer] {
	return func(ctx *gin.Context, rows *[]Customer) error {
		for i := range *rows {
			if (*rows)[i].Img != "" {
				(*rows)[i].Img = filesystem.FullImageURL(ctx, (*rows)[i].Img)
			}
		}
		return nil
	}
}

func beforeResponse() basehandler.TResponseHook[Customer] {
	return func(ctx *gin.Context, cus *Customer) error {
		if cus.Img != "" {
			cus.Img = filesystem.FullImageURL(ctx, cus.Img)
		}
		return nil
	}
}
