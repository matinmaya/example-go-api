package authctx

import (
	"context"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func SetUserID(ctx context.Context, userID uint32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) *uint32 {
	val := ctx.Value(userIDKey)
	switch v := val.(type) {
	case uint32:
		return &v
	case int:
		converted := uint32(v)
		return &converted
	}
	return nil
}
