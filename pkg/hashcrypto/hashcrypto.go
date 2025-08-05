package hashcrypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/spaolacci/murmur3"
	"golang.org/x/crypto/bcrypt"
)

func HashMake(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func HashCheck(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashMakeToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func HashCheckToken(token string, hashed string) bool {
	return HashMakeToken(token) == hashed
}

func HashCacheKey(s string) string {
	h := murmur3.Sum64([]byte(s))
	return fmt.Sprintf("%x", h)
}
