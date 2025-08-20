package module

import (
	"github.com/gin-gonic/gin"

	"reapp/pkg/base/basehandler"
)

type HandlerYields[T any, TDTO any] struct {
	CreateValidateScope basehandler.TScope[TDTO]
	UpdateValidateScope basehandler.TScopeWithID[TDTO]
	AfterValidate       basehandler.TAfterValidate[TDTO]
	BeforeResponse      basehandler.TBeforeResponse[T]
	BeforeResponseList  basehandler.TBeforeResponseList[T]
}

type Handler[T IWithID[TID], TID IUintID, TDTO any, TQ any] struct {
	service IServiceAdapter[T, TID]
	yields  HandlerYields[T, TDTO]
}

func NewHandler[T IWithID[TID], TID IUintID, TDTO any, TQ any](
	service IServiceAdapter[T, TID],
	yields HandlerYields[T, TDTO],
) IHandler[T, TID, TDTO, TQ] {
	return Handler[T, TID, TDTO, TQ]{service, yields}
}

func (h Handler[T, TID, TDTO, TQ]) Create(ctx *gin.Context) {
	var model T
	var modelDTO TDTO

	basehandler.Create(ctx,
		h.service,
		&model,
		&modelDTO,
		h.yields.CreateValidateScope,
		h.yields.AfterValidate,
		h.yields.BeforeResponse,
	)
}

func (h Handler[T, TID, TDTO, TQ]) Update(ctx *gin.Context) {
	var modelDTO TDTO
	basehandler.Update(ctx,
		h.service,
		&modelDTO,
		h.yields.UpdateValidateScope,
		h.yields.AfterValidate,
		h.yields.BeforeResponse,
	)
}

func (h Handler[T, TID, TDTO, TQ]) Delete(ctx *gin.Context) {
	basehandler.Delete(ctx, h.service, nil)
}

func (h Handler[T, TID, TDTO, TQ]) List(ctx *gin.Context) {
	var query TQ
	basehandler.Paginate(ctx,
		h.service,
		&query,
		h.yields.BeforeResponseList,
	)
}
