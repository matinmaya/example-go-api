package reqctx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func SnakeToCamelCase(snake string) string {
	if snake == "" {
		return ""
	}
	if strings.ToLower(snake) == "id" {
		return "ID"
	}

	parts := strings.Split(snake, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(parts, "")
}

func RestoreRequestBody(ctx *gin.Context, raw []byte) {
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
}

func GetFieldNames(ctx *gin.Context) ([]string, error) {
	var fieldNames []string
	if ctx == nil || ctx.Request == nil || ctx.Request.Body == nil {
		return fieldNames, nil
	}

	raw, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	RestoreRequestBody(ctx, raw)

	var fields map[string]interface{}
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	guarded := map[string]struct{}{
		"id": {}, "created_by": {}, "updated_by": {}, "deleted_by": {},
		"created_at": {}, "updated_at": {}, "deleted_at": {},
	}

	for key := range fields {
		if key == "" {
			continue
		}
		if _, isGuarded := guarded[key]; isGuarded {
			continue
		}
		fieldNames = append(fieldNames, SnakeToCamelCase(key))
	}

	return fieldNames, nil
}

func RemoveFields(fields *[]string, targets ...string) {
	targetSet := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		targetSet[strings.ToLower(t)] = struct{}{}
	}

	keptFields := (*fields)[:0]
	for _, field := range *fields {
		if _, found := targetSet[strings.ToLower(field)]; !found {
			keptFields = append(keptFields, field)
		}
	}
	*fields = keptFields
}
