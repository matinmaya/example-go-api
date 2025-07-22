package jwthelper

import (
	"reapp/internal/modules/user/usermodel"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var JWTSecret = []byte("")
var AccessTokenTTL = 60 * time.Minute
var RefreshTokenTTL = 30 * 24 * time.Hour

type Claims struct {
	UserID   uint32   `json:"user_id"`
	RoleIDs  []uint16 `json:"role_ids"`
	Username string   `json:"username"`
	jwt.StandardClaims
}

// InitJWT initializes JWT helper settings including the secret key and token lifetimes.
// It sets the signing secret and defines the expiration durations for access and refresh tokens.
// Parameters:
// - secret: the secret key used to sign JWT tokens.
// - accessTokenTTL: duration in minutes before the access token expires.
// - refreshTokenTTL: duration in minutes before the refresh token expires.
func InitJWT(secret string, accessTokenTTL int, refreshTokenTTL int) {
	JWTSecret = []byte(secret)
	AccessTokenTTL = time.Duration(accessTokenTTL) * time.Minute
	RefreshTokenTTL = time.Duration(refreshTokenTTL) * time.Minute
}

func ExistsSecret() bool {
	return len(JWTSecret) > 0
}

func GenerateTokenWithExpiry(user usermodel.User, duration time.Duration) (string, jwt.StandardClaims, error) {
	expirationTime := time.Now().Add(duration)
	var roleIDs []uint16
	for _, r := range user.Roles {
		roleIDs = append(roleIDs, uint16(r.ID))
	}

	jti := uuid.New().String()
	standardClaims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        jti,
	}

	claims := &Claims{
		UserID:         user.ID,
		RoleIDs:        roleIDs,
		Username:       user.Username,
		StandardClaims: standardClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	generatedToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", jwt.StandardClaims{}, err
	}
	return generatedToken, standardClaims, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
