package userhandler

import (
	"fmt"
	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/base/basehandler"
	"reapp/pkg/context/dbctx"
	"reapp/pkg/http/reqctx"
	"reapp/pkg/http/reqvalidate"
	"reapp/pkg/http/response"
	"reapp/pkg/services/rediservice"
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
	basehandler.Create(ctx, h.service, &usermodel.User{}, &usermodel.User{}, nil, func(user *usermodel.User) error {
		user.RoleIds = []uint16{}
		return nil
	})
}

func (h *UserHandler) Update(ctx *gin.Context) {
	basehandler.Update(ctx, h.service, &usermodel.User{}, func(user *usermodel.User, id uint64) error {
		user.UniqueScope = validators.ExceptByID(id)
		return nil
	}, func(fields *[]string) error {
		reqctx.RemoveFields(fields, "Password")
		return nil
	}, func(user *usermodel.User) error {
		user.RoleIds = []uint16{}
		go rediservice.RemoveCacheOfAuthUser(fmt.Sprintf("%v", user.ID))
		return nil
	})
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	basehandler.Delete(ctx, h.service, func(ctx *gin.Context) error {
		userID, _ := ctx.Get("user_id")
		go rediservice.RemoveCacheOfAuthUser(fmt.Sprintf("%v", userID))
		return nil
	})
}

func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	db := dbctx.DB(ctx)
	var data usermodel.ChangePassword
	if !reqvalidate.Validate(ctx, &data) {
		return
	}

	err := h.service.ChangePassword(db, data)
	response.AsJSON(ctx, nil, err)
}
