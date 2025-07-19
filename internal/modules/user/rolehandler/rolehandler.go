package rolehandler

import (
	"net/http"
	"reapp/internal/handler"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/roleservice"
	"reapp/pkg/binding"
	"reapp/pkg/mapper"
	"reapp/pkg/requestutils"
	"reapp/pkg/response"
	"reapp/pkg/validators"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service roleservice.IRoleService
}

func NewRoleHandler(s roleservice.IRoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

func (h *RoleHandler) Create(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	fields, bad := requestutils.GetFieldNames(ctx)
	if bad != nil {
		response.Error(ctx, http.StatusBadRequest, bad.Error(), nil)
		return
	}

	var role rolemodel.Role
	var roleDto rolemodel.Role
	if !binding.ValidateData(ctx, &roleDto) {
		return
	}

	if err := mapper.MapStruct(&role, roleDto, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.service.Create(db, &role)
	response.AsJSON(ctx, role, err)
}

func (h *RoleHandler) Update(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	fields, bad := requestutils.GetFieldNames(ctx)
	if bad != nil {
		response.Error(ctx, http.StatusBadRequest, bad.Error(), nil)
		return
	}

	role, fail := h.service.GetByID(db, uint16(id))
	if fail != nil {
		response.Error(ctx, http.StatusBadRequest, fail.Error(), nil)
		return
	}

	var roleDto rolemodel.Role
	roleDto.ScopeUnique = validators.ExceptByID(uint64(id))
	if !binding.ValidateData(ctx, &roleDto) {
		return
	}

	if err := mapper.MapStruct(role, roleDto, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.service.Update(db, role)
	response.AsJSON(ctx, role, err)
}

func (h *RoleHandler) GetDetail(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	role, err := h.service.GetDetail(db, uint16(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "success", role)
}

func (h *RoleHandler) Delete(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	err := h.service.Delete(db, uint16(id))
	response.AsJSON(ctx, nil, err)
}

func (h *RoleHandler) GetAll(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	data, err := h.service.GetAll(db)
	response.AsJSON(ctx, data, err)
}

func (h *RoleHandler) List(ctx *gin.Context) {
	var query rolemodel.RoleListQuery
	handler.PaginateList(ctx, &query, h.service)
}
