package userhandler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"reapp/internal/modules/user/usermodel"
	"reapp/internal/modules/user/userservice"
	"reapp/pkg/base/basehandler"
	"reapp/pkg/context/dbctx"
	"reapp/pkg/filesystem"
	"reapp/pkg/http/reqctx"
	"reapp/pkg/http/reqvalidate"
	"reapp/pkg/http/response"
	"reapp/pkg/services/rediservice"
	"reapp/pkg/validators"
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
	basehandler.Create(ctx, h.service, &usermodel.User{}, &usermodel.User{}, nil, afterValidate(false), beforeResponse(false))
}

func (h *UserHandler) Update(ctx *gin.Context) {
	basehandler.Update(ctx, h.service, &usermodel.User{}, func(user *usermodel.User, id uint64) error {
		user.UniqueScope = validators.ExceptByID(id)
		return nil
	}, afterValidate(true), beforeResponse(true))
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

func afterValidate(isUpdate bool) func(ctx *gin.Context, user *usermodel.User, fields *[]string) error {
	return func(ctx *gin.Context, user *usermodel.User, fields *[]string) error {
		if isUpdate {
			reqctx.RemoveFields(fields, "Password")
		}
		if user.Img != "" {
			if filesystem.IsFullImagePath(user.Img) {
				reqctx.RemoveFields(fields, "Img")
			} else if !filesystem.IsAbsoluteImagePath(user.Img) {
				return errors.New("invalid image path format: must start with '/' or be a full URL")
			}
		}
		return nil
	}
}

func beforeResponse(isUpdate bool) func(ctx *gin.Context, user *usermodel.User) error {
	return func(ctx *gin.Context, user *usermodel.User) error {
		user.RoleIds = []uint16{}
		if user.Img != "" {
			user.Img = filesystem.FullImageURL(ctx, user.Img)
		}
		if isUpdate {
			go rediservice.RemoveCacheOfAuthUser(fmt.Sprintf("%v", user.ID))
		}
		return nil
	}
}
