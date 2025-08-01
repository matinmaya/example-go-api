package userhandler

import (
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/base/basehandler"
	"reapp/pkg/binding"
	"reapp/pkg/helpers/ctxhelper"
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

func (h *UserHandler) List(ctx *gin.Context) {
	basehandler.Paginate(ctx, h.service, &usermodel.UserListQuery{})
}

func (h *UserHandler) Create(ctx *gin.Context) {
	basehandler.Create(ctx, h.service, &usermodel.User{}, &usermodel.User{}, nil, nil)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	basehandler.Update(ctx, h.service, &usermodel.User{}, func(modelDTO any, id uint64) error {
		if dto, ok := modelDTO.(*usermodel.User); ok {
			dto.ScopeUnique = validators.ExceptByID(id)
		}
		return nil
	}, func(fields *[]string) error {
		requestutils.RemoveFields(fields, "Password")
		return nil
	}, nil)
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	basehandler.Delete(ctx, h.service)
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
