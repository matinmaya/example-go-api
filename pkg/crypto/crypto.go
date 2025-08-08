package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/spaolacci/murmur3"
	"golang.org/x/crypto/bcrypt"
)

func Make(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func Check(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MakeToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func CheckToken(token string, hashed string) bool {
	return MakeToken(token) == hashed
}

func CacheKey(s string) string {
	h := murmur3.Sum64([]byte(s))
	return fmt.Sprintf("%x", h)
}
