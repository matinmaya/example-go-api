package handler

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/service"
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"
	"reapp/pkg/response"
	"reflect"

	"github.com/gin-gonic/gin"
)

func PaginateList(ctx *gin.Context, query any, moduleService service.Lister) {
	db := ctxhelper.GetDB(ctx)
	var pagination paginator.Pagination

	valueOfQuery := reflect.ValueOf(query)
	if valueOfQuery.Kind() != reflect.Ptr {
		response.Error(ctx, http.StatusInternalServerError, "the destination data must be provided as a reference", nil)
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

	if err := moduleService.List(db, &pagination, filters); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.AsJSON(ctx, pagination, nil)
}
