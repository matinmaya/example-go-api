package authhandler

import (
	"net/http"
	"reapp/internal/helpers/ctxhelper"
	"reapp/internal/helpers/jwthelper"
	"reapp/internal/helpers/loghelper"
	"reapp/internal/modules/user/authservice"
	"reapp/internal/modules/user/usermodel"
	"reapp/pkg/basemodel"
	"reapp/pkg/binding"
	"reapp/pkg/hashcrypto"
	"reapp/pkg/response"
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

	accessToken, _, fail := jwthelper.GenerateTokenWithExpiry(*user, jwthelper.AccessTokenTTL)
	if fail != nil {
		response.Error(ctx, http.StatusUnauthorized, fail.Error(), nil)
		return
	}

	refreshToken, refreshClaims, fail := jwthelper.GenerateTokenWithExpiry(*user, jwthelper.RefreshTokenTTL)
	if fail != nil {
		response.Error(ctx, http.StatusUnauthorized, fail.Error(), nil)
		return
	}

	hashedRefreshToken := hashcrypto.HashMakeToken(refreshToken)
	if hashedRefreshToken == "" {
		response.Error(ctx, http.StatusInternalServerError, "Failed to hash refresh token", nil)
		return
	}

	uaString := ctx.Request.UserAgent()
	ip := ctx.ClientIP()
	device, platform, browser, os := loghelper.ParseUserAgent(uaString)
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
		response.Error(ctx, http.StatusInternalServerError, "Failed to save token info", nil)
		return
	}

	data := &usermodel.AuthLoginResource{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.Success(ctx, http.StatusOK, "login success", data)
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	db := ctxhelper.GetDB(ctx)
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if !binding.ValidateData(ctx, &req) {
		return
	}

	token, isFound := ctx.Get("jwt_token")
	if !isFound {
		response.Error(ctx, http.StatusUnauthorized, "No token found", nil)
		return
	}

	claims, ok := token.(*jwthelper.Claims)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Invalid token claims", nil)
		return
	}

	odlRefreshClaims, err := jwthelper.ParseToken(req.RefreshToken)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	userID := odlRefreshClaims.UserID
	jti := odlRefreshClaims.Id
	tokenInfo, err := h.service.GetTokenInfo(db, userID, jti)
	if err != nil || tokenInfo == nil {
		response.Error(ctx, http.StatusUnauthorized, "Refresh token not found", nil)
		return
	}

	if time.Now().After(tokenInfo.ExpiresAt.Time) {
		response.Error(ctx, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	if !hashcrypto.HashCheckToken(req.RefreshToken, tokenInfo.RefreshToken) {
		response.Error(ctx, http.StatusUnauthorized, "Refresh token mismatch", nil)
		return
	}

	user, err := h.service.GetUserByID(db, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "User not found", nil)
		return
	}

	accessToken, _, err := jwthelper.GenerateTokenWithExpiry(*user, jwthelper.AccessTokenTTL)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to generate access token", nil)
		return
	}

	newRefreshToken, refreshClaims, err := jwthelper.GenerateTokenWithExpiry(*user, jwthelper.RefreshTokenTTL)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to generate refresh token", nil)
		return
	}

	if err := h.service.RevokeAccessToken(claims.Id, time.Unix(claims.ExpiresAt, 0)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to revoke old access token", nil)
		return
	}

	hashedNewRefresh := hashcrypto.HashMakeToken(newRefreshToken)
	if hashedNewRefresh == "" {
		response.Error(ctx, http.StatusInternalServerError, "Failed to hash refresh token", nil)
		return
	}

	tokenInfo.JTI = refreshClaims.Id
	tokenInfo.RefreshToken = hashedNewRefresh
	tokenInfo.ExpiresAt = basemodel.DateTimeFormat{Time: time.Unix(refreshClaims.ExpiresAt, 0)}
	if err := h.service.UpdateTokenInfo(db, tokenInfo); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to update token info", nil)
		return
	}

	data := &usermodel.AuthLoginResource{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}
	response.Success(ctx, http.StatusOK, "token refreshed", data)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	token, exists := ctx.Get("jwt_token")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "No token found", nil)
		return
	}

	claims, ok := token.(*jwthelper.Claims)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Invalid token claims", nil)
		return
	}

	err := h.service.RevokeAccessToken(claims.Id, time.Unix(claims.ExpiresAt, 0))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to revoke token", nil)
		return
	}

	db := ctxhelper.GetDB(ctx)
	if err := h.service.DeleteTokenInfoByUserID(db, claims.UserID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete token info", nil)
		return
	}

	response.Success(ctx, http.StatusOK, "logout success", nil)
}
