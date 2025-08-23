package basemodel

import (
	"context"
	"reflect"

	"gorm.io/gorm"
)

type ctxKey string

const oldValueKey ctxKey = "oldValue"

func SetOldValue(tx *gorm.DB, instance any) error {
	data := make(map[string]interface{})

	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, field := range tx.Statement.Schema.Fields {
		fieldVal, _ := field.ValueOf(tx.Statement.Context, val)
		data[field.Name] = fieldVal
	}

	ctx := context.WithValue(tx.Statement.Context, oldValueKey, data)
	tx.Statement.Context = ctx
	return nil
}

func OldValue(tx *gorm.DB) map[string]interface{} {
	data := tx.Statement.Context.Value(oldValueKey)
	if data == nil {
		return nil
	}
	return data.(map[string]interface{})
}
