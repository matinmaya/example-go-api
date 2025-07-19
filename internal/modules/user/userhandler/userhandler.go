package userhandler

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/binding"
	"reapp/pkg/mapper"
	"reapp/pkg/requestutils"
	"reapp/pkg/response"
	"reapp/pkg/validators"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service userservice.IUserService
}

func NewUserHandler(s userservice.IUserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	fields, bad := requestutils.GetFieldNames(ctx)
	if bad != nil {
		response.Error(ctx, http.StatusBadRequest, bad.Error(), nil)
		return
	}

	var user usermodel.User
	var userDto usermodel.User
	if !binding.ValidateData(ctx, &userDto) {
		return
	}

	if err := mapper.MapStruct(&user, userDto, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.service.Create(db, &user)
	response.AsJSON(ctx, nil, err)
}

func (h *UserHandler) Update(ctx *gin.Context) {
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

	requestutils.RemoveFields(&fields, "password")
	user, fail := h.service.GetByID(db, uint32(id))
	if fail != nil {
		response.Error(ctx, http.StatusNotFound, fail.Error(), nil)
		return
	}

	var userDto usermodel.User
	userDto.ScopeUnique = validators.ExceptByID(uint64(id))
	if !binding.ValidateData(ctx, &userDto) {
		return
	}

	if err := mapper.MapStruct(user, userDto, fields); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := h.service.Update(db, user)
	response.AsJSON(ctx, user, err)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var id uint64
	if !binding.ValidateParamID(ctx, &id) {
		return
	}

	err := h.service.Delete(db, uint32(id))
	response.AsJSON(ctx, nil, err)
}

func (h *UserHandler) GetAll(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	data, err := h.service.GetAll(db)
	response.AsJSON(ctx, data, err)
}

func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var data usermodel.ChangePassword
	if !binding.ValidateData(ctx, &data) {
		return
	}

	err := h.service.ChangePassword(db, data)
	response.AsJSON(ctx, nil, err)
}
