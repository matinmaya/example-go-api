package jwthelper

import (
	"reapp/internal/modules/user/usermodel"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTSecret = []byte("")

type Claims struct {
	UserID   uint32   `json:"user_id"`
	RoleIDs  []uint16 `json:"role_ids"`
	Username string   `json:"username"`
	jwt.StandardClaims
}

func SetSecret(secret string) {
	JWTSecret = []byte(secret)
}

func GetSecret() []byte {
	return JWTSecret
}

func GenerateJWT(user usermodel.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	var roleIDs []uint16
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, uint16(r.ID))
	}

	claims := &Claims{
		UserID:   user.ID,
		RoleIDs:  roleIDs,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
