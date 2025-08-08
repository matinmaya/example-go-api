package authhandler

import (
	"net/http"
	"reapp/internal/modules/user/authservice"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/base/basemodel"
	"reapp/pkg/context/dbctx"
	"reapp/pkg/crypto"
	"reapp/pkg/http/reqvalidate"
	"reapp/pkg/http/response"
	"reapp/pkg/lang"
	"reapp/pkg/logger"
	"reapp/pkg/services/jwtservice"
	"reapp/pkg/services/rediservice"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service authservice.IAuthService
}

func NewAuthHandler(s authservice.IAuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	db := dbctx.DB(ctx)
	var credentials usermodel.AuthCredentials
	if !reqvalidate.Validate(ctx, &credentials) {
		return
	}

	user, err := h.service.Attempt(db, credentials)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	if !user.Status {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "account_locked"), nil)
		return
	}

	accessToken, accessClaims, fail := jwtservice.GenerateTokenWithExpiry(*user, jwtservice.AccessTokenTTL)
	if fail != nil {
		response.Error(ctx, http.StatusUnauthorized, fail.Error(), nil)
		return
	}

	refreshToken, refreshClaims, fail := jwtservice.GenerateTokenWithExpiry(*user, jwtservice.RefreshTokenTTL)
	if fail != nil {
		response.Error(ctx, http.StatusUnauthorized, fail.Error(), nil)
		return
	}

	hashedRefreshToken := crypto.MakeToken(refreshToken)
	if hashedRefreshToken == "" {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_to_hash_refresh"), nil)
		return
	}

	uaString := ctx.Request.UserAgent()
	ip := ctx.ClientIP()
	device, platform, browser, os := logger.ParseUserAgent(uaString)
	expiresAt := basemodel.DateTimeFormat{Time: time.Unix(refreshClaims.ExpiresAt, 0)}
	tokenInfo := &usermodel.TokenInfo{
		UserID:       user.ID,
		JTI:          refreshClaims.Id,
		RefreshToken: hashedRefreshToken,
		Device:       device,
		Platform:     platform,
		Browser:      browser,
		OS:           os,
		UserAgent:    uaString,
		IP:           ip,
		ExpiresAt:    expiresAt,
	}
	if err := h.service.SaveTokenInfo(db, tokenInfo); err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_to_save_token"), nil)
		return
	}

	data := &usermodel.AuthLoginResource{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    basemodel.DateTimeFormat{Time: time.Unix(accessClaims.ExpiresAt, 0)},
	}
	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "auth", "login_success"), data)
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	db := dbctx.DB(ctx)
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if !reqvalidate.Validate(ctx, &req) {
		return
	}

	odlRefreshClaims, err := jwtservice.ParseToken(req.RefreshToken)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "invalid_refresh_token"), nil)
		return
	}

	userID := odlRefreshClaims.UserID
	jti := odlRefreshClaims.Id
	tokenInfo, err := h.service.GetTokenInfo(db, userID, jti)
	if err != nil || tokenInfo == nil {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "refresh_token_not_found"), nil)
		return
	}

	if time.Now().After(tokenInfo.ExpiresAt.Time) {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "refresh_token_expired"), nil)
		return
	}

	if !crypto.CheckToken(req.RefreshToken, tokenInfo.RefreshToken) {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "refresh_token_mismatch"), nil)
		return
	}

	user, err := h.service.GetUserByID(db, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "user_not_found_internal"), nil)
		return
	}

	if !user.Status {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "account_locked"), nil)
		return
	}

	accessToken, accessClaims, err := jwtservice.GenerateTokenWithExpiry(*user, jwtservice.AccessTokenTTL)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_generate_access"), nil)
		return
	}

	newRefreshToken, refreshClaims, err := jwtservice.GenerateTokenWithExpiry(*user, jwtservice.RefreshTokenTTL)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_generate_refresh"), nil)
		return
	}

	if token, exists := ctx.Get("jwt_token"); exists {
		if claims, ok := token.(*jwtservice.Claims); ok {
			if err := rediservice.RevokeToken(claims.Id, time.Unix(claims.ExpiresAt, 0)); err != nil {
				response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_revoke_access"), nil)
				return
			}
		}
	}

	hashedNewRefresh := crypto.MakeToken(newRefreshToken)
	if hashedNewRefresh == "" {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_to_hash_refresh"), nil)
		return
	}

	tokenInfo.JTI = refreshClaims.Id
	tokenInfo.RefreshToken = hashedNewRefresh
	tokenInfo.ExpiresAt = basemodel.DateTimeFormat{Time: time.Unix(refreshClaims.ExpiresAt, 0)}
	if err := h.service.UpdateTokenInfo(db, tokenInfo); err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_update_token"), nil)
		return
	}

	data := &usermodel.AuthLoginResource{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    basemodel.DateTimeFormat{Time: time.Unix(accessClaims.ExpiresAt, 0)},
	}
	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "auth", "token_refreshed"), data)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	token, exists := ctx.Get("jwt_token")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "no_token_found"), nil)
		return
	}

	claims, ok := token.(*jwtservice.Claims)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, lang.Tran(ctx, "auth", "invalid_token_claims"), nil)
		return
	}

	err := rediservice.RevokeToken(claims.Id, time.Unix(claims.ExpiresAt, 0))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_revoke_token"), nil)
		return
	}

	db := dbctx.DB(ctx)
	if err := h.service.DeleteTokenInfo(db, claims.UserID, claims.Id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, lang.Tran(ctx, "auth", "failed_delete_token"), nil)
		return
	}

	response.Success(ctx, http.StatusOK, lang.Tran(ctx, "auth", "logout_success"), nil)
}
