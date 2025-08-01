package basehandler

import (
	"fmt"
	"net/http"
	"reapp/pkg/binding"
	"reapp/pkg/filterscopes"
	"reapp/pkg/helpers/ctxhelper"
	"reapp/pkg/lang"
	"reapp/pkg/mapper"
	"reapp/pkg/paginator"
	"reapp/pkg/requestutils"
	"reapp/pkg/response"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TScopeFunc func(any) error
type TScopeWithIDFunc func(any, uint64) error

type IServiceLister interface {
	List(db *gorm.DB, pagination *paginator.Pagination, filters []filterscopes.QueryFilter) error
}

type IServiceGetter[T any] interface {
	GetAll(db *gorm.DB) ([]T, error)
}

type IServiceGetterDetail[T any] interface {
	GetDetail(db *gorm.DB, id uint64) (*T, error)
}

type IServiceCreator[T any] interface {
	Create(db *gorm.DB, model T) error
}

type IServiceUpdater[T any] interface {
	Update(db *gorm.DB, model *T) error
	GetByID(db *gorm.DB, id uint64) (*T, error)
}

type IServiceDeleter interface {
	Delete(db *gorm.DB, id uint64) error
}

func Paginate(ctx *gin.Context, service IServiceLister, query any) {
	db := ctxhelper.GetDB(ctx)
	var pagination paginator.Pagination

	valueOfQuery := reflect.ValueOf(query)
	if valueOfQuery.Kind() != reflect.Ptr {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "internal", "required_pointer"), nil)
		return
	}

	if err := ctx.ShouldBindQuery(query); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	queryValues := ctx.Request.URL.Query()
	filters := filterscopes.ParseQueryByUrlValues(query, queryValues)

	if err := service.List(db, &pagination, filters); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.AsJSON(ctx, pagination, nil)
}

func GetAll[T any](ctx *gin.Context, service IServiceGetter[T]) {
	db := ctxhelper.GetDB(ctx)

	data, err := service.GetAll(db)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.AsJSON(ctx, data, nil)
}

func GetDetail[T any](ctx *gin.Context, service IServiceGetterDetail[T]) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	data, err := service.GetDetail(db, id)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, lang.Tran(ctx, "response", "not_found"), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), data)
}

func Create[T any](
	ctx *gin.Context,
	service IServiceCreator[T],
	model T,
	modelDTO any,
	setValidationScope TScopeFunc,
) {
	db := ctxhelper.GetDB(ctx)
	valueOfModelDTO := reflect.ValueOf(modelDTO)
	if valueOfModelDTO.Kind() != reflect.Ptr {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "internal", "required_pointer"), nil)
		return
	}

	fields, bad := requestutils.GetFieldNames(ctx)
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

	if !binding.ValidateData(ctx, modelDTO) {
		return
	}

	if err := mapper.MapStruct(model, modelDTO, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := service.Create(db, model)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), model)
}

func Update[T any](
	ctx *gin.Context,
	service IServiceUpdater[T],
	modelDTO any,
	setValidationScope TScopeWithIDFunc,
	removeFields TScopeFunc,
) {
	db := ctxhelper.GetDB(ctx)
	valueOfModelDTO := reflect.ValueOf(modelDTO)
	if valueOfModelDTO.Kind() != reflect.Ptr {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "internal", "required_pointer"), nil)
		return
	}

	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	fields, bad := requestutils.GetFieldNames(ctx)
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
	if !binding.ValidateData(ctx, modelDTO) {
		return
	}

	if removeFields != nil {
		if err := removeFields(&fields); err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}
	fmt.Printf("Fields to update: %v\n", fields)
	if err := mapper.MapStruct(model, modelDTO, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := service.Update(db, model)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), model)
}

func Delete(ctx *gin.Context, service IServiceDeleter) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	err := service.Delete(db, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "response", "success"), nil)
}
