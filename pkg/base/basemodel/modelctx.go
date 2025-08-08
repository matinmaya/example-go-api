package basemodel

import (
	"context"
	"reflect"

	"gorm.io/gorm"
)

type ctxKey string

const bfChangeValueKey ctxKey = "bfChangeValue"

func SetBfChangeValueFromInstance(db *gorm.DB, instance any) error {
	data := make(map[string]interface{})

	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, field := range db.Statement.Schema.Fields {
		fieldVal, _ := field.ValueOf(db.Statement.Context, val)
		data[field.Name] = fieldVal
	}

	ctx := context.WithValue(db.Statement.Context, bfChangeValueKey, data)
	db.Statement.Context = ctx
	return nil
}

func GetBfChangeValue(tx *gorm.DB) map[string]interface{} {
	data := tx.Statement.Context.Value(bfChangeValueKey)
	if data == nil {
		return nil
	}
	return data.(map[string]interface{})
}
