package basehandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"reapp/pkg/context/dbctx"
	"reapp/pkg/http/reqctx"
	"reapp/pkg/http/reqvalidate"
	"reapp/pkg/http/response"
	"reapp/pkg/lang"
	"reapp/pkg/mapper"
	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
)

func Paginate[T any, TQ any](
	ctx *gin.Context,
	service IServiceLister[T],
	query *TQ,
	beforeResponse TBeforeResponseList[T],
) {
	db := dbctx.DB(ctx)
	var pagination paginator.Pagination[T]

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Printf("%s", err.Error())
		response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "invalid_query_params"), nil)
		return
	}

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		log.Printf("%s", err.Error())
		response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "invalid_query_params"), nil)
		return
	}

	values := ctx.Request.URL.Query()
	filterFields := queryfilter.FilterFields(query, values)
	if err := service.List(ctx, db, &pagination, filterFields); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if beforeResponse != nil {
		if err := beforeResponse(ctx, &pagination.Rows); err != nil {
			response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	response.AsJSON(ctx, pagination, nil)
}

func GetAll[T any](ctx *gin.Context, service IServiceGetter[T]) {
	db := dbctx.DB(ctx)

	data, err := service.GetAll(db)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.AsJSON(ctx, data, nil)
}

func GetDetail[T any](ctx *gin.Context, service IServiceGetterDetail[T]) {
	db := dbctx.DB(ctx)
	var id uint64
	if !reqvalidate.ValidateParamID(ctx, &id) {
		return
	}

	data, err := service.GetDetail(db, id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, lang.Tran(ctx, "response", "not_found"), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), data)
}

func Create[T1 any, T2 any](
	ctx *gin.Context,
	service IServiceCreator[T1],
	model *T1,
	modelDTO *T2,
	setValidationScope TScope[T2],
	afterValidate TAfterValidate[T2],
	beforeResponse TBeforeResponse[T1],
) {
	db := dbctx.DB(ctx)
	fields, bad := reqctx.GetFieldNames(ctx)
	if bad != nil {
		response.Error(ctx, http.StatusBadRequest, bad.Error(), nil)
		return
	}

	if setValidationScope != nil {
		if err := setValidationScope(modelDTO); err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	if !reqvalidate.Validate(ctx, modelDTO) {
		return
	}

	if afterValidate != nil {
		if err := afterValidate(ctx, modelDTO, &fields); err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	if err := mapper.MapModel(model, modelDTO, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := service.Create(db, model)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if beforeResponse != nil {
		if err := beforeResponse(ctx, model); err != nil {
			response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), model)
}

func Update[T1 any, T2 any](
	ctx *gin.Context,
	service IServiceUpdater[T1],
	modelDTO *T2,
	setValidationScope TScopeWithID[T2],
	afterValidate TAfterValidate[T2],
	beforeResponse TBeforeResponse[T1],
) {
	db := dbctx.DB(ctx)
	var id uint64
	if !reqvalidate.ValidateParamID(ctx, &id) {
		return
	}

	fields, bad := reqctx.GetFieldNames(ctx)
	if bad != nil {
		response.Error(ctx, http.StatusBadRequest, bad.Error(), nil)
		return
	}

	model, fail := service.GetByID(db, id)
	if fail != nil {
		response.Error(ctx, http.StatusNotFound, lang.Tran(ctx, "response", "not_found"), nil)
		return
	}

	if setValidationScope != nil {
		if err := setValidationScope(modelDTO, id); err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}
	if !reqvalidate.Validate(ctx, modelDTO) {
		return
	}

	if afterValidate != nil {
		if err := afterValidate(ctx, modelDTO, &fields); err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}

	if err := mapper.MapModel(model, modelDTO, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := service.Update(db, model)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if beforeResponse != nil {
		if err := beforeResponse(ctx, model); err != nil {
			response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), model)
}

func Delete(ctx *gin.Context, service IServiceDeleter, beforeResponse func(*gin.Context) error) {
	db := dbctx.DB(ctx)
	var id uint64
	if !reqvalidate.ValidateParamID(ctx, &id) {
		return
	}

	err := service.Delete(db, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if beforeResponse != nil {
		if err := beforeResponse(ctx); err != nil {
			response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), nil)
}
