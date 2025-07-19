package authhandler

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/helpers/jwthelper"
	"reapp/internal/modules/user/authservice"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/binding"
	"reapp/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service authservice.IAuthService
}

func NewAuthHandler(s authservice.IAuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var credentials usermodel.AuthCredentials
	if !binding.ValidateData(ctx, &credentials) {
		return
	}

	user, err := h.service.Attempt(db, credentials)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	token, fail := jwthelper.GenerateJWT(*user)
	if fail != nil {
		response.Error(ctx, http.StatusUnauthorized, fail.Error(), nil)
		return
	}

	data := &usermodel.AuthLoginResource{
		Username: user.Username,
		Token:    token,
	}
	response.Success(ctx, http.StatusOK, "Login Success", data)
}
