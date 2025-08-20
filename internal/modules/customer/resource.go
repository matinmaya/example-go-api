package customer

import (
	"github.com/gin-gonic/gin"

	"reapp/pkg/filesystem"
	"reapp/pkg/validators"
)

func updateValidateScope() func(*Customer, uint64) error {
	return func(cus *Customer, id uint64) error {
		cus.UniqueScope = validators.ExceptByID(id)
		return nil
	}
}

func listBeforeResponse() func(*gin.Context, *[]Customer) error {
	return func(ctx *gin.Context, rows *[]Customer) error {
		for i := range *rows {
			if (*rows)[i].Img != "" {
				(*rows)[i].Img = filesystem.FullImageURL(ctx, (*rows)[i].Img)
			}
		}
		return nil
	}
}

func beforeResponse() func(*gin.Context, *Customer) error {
	return func(ctx *gin.Context, cus *Customer) error {
		if cus.Img != "" {
			cus.Img = filesystem.FullImageURL(ctx, cus.Img)
		}
		return nil
	}
}
