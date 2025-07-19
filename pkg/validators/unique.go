package validators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ValidateScopeUnique struct {
	ScopeUnique func() func(db *gorm.DB) *gorm.DB `json:"-" gorm:"-"`
}

func Unique(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	params := strings.Split(fl.Param(), "?")
	if len(params) < 1 {
		fmt.Printf("Messing table name in unique validation")
		return false
	}

	tableName := params[0]
	// uniqueFields := make(map[string]interface{})

	structFieldName := fl.StructFieldName()

	field, ok := reflect.TypeOf(fl.Parent().Interface()).FieldByName(structFieldName)
	if !ok {
		fmt.Printf("Struct field %s not found in struct", structFieldName)
		return false
	}

	jsonTag := field.Tag.Get("json")
	columnName := strings.Split(jsonTag, ",")[0]

	var count int64
	q := DB.Table(tableName).Where(columnName+" = ?", fieldValue)

	scopeMethod := reflect.ValueOf(fl.Parent().Interface()).FieldByName("ScopeUnique")
	if scopeMethod.IsValid() && !scopeMethod.IsNil() {
		results := scopeMethod.Call(nil)
		if len(results) > 0 {
			if scopeFunc, ok := results[0].Interface().(func(*gorm.DB) *gorm.DB); ok {
				q.Scopes(scopeFunc)
			}
		}
	}

	q.Count(&count)

	return count == 0
}

func ExceptByID(id uint64, prefix ...string) func() func(*gorm.DB) *gorm.DB {
	return func() func(tx *gorm.DB) *gorm.DB {
		return func(tx *gorm.DB) *gorm.DB {
			if len(prefix) > 0 && prefix[0] != "" {
				return tx.Where(fmt.Sprintf("%s.id = ?", prefix[0]), id)
			}
			return tx.Where("id <> ?", id)
		}
	}
}
