package queryfilter

import (
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type FilterField struct {
	Column   string
	Operator string
	Value    string
}

func FilterFields[T any](queryDTO T, values url.Values) []FilterField {
	var filterFields []FilterField
	valOfQuery := reflect.ValueOf(queryDTO)
	if valOfQuery.Kind() == reflect.Ptr {
		valOfQuery = valOfQuery.Elem()
	}

	typeOfQuery := valOfQuery.Type()
	for i := 0; i < valOfQuery.NumField(); i++ {
		structField := typeOfQuery.Field(i)
		filterTag := structField.Tag.Get("filter")
		if filterTag == "" {
			continue
		}

		param := structField.Tag.Get("form")
		if param == "" {
			param = structField.Name
		}
		if val := values.Get(param); val == "" {
			continue
		}

		value := fieldValue(valOfQuery.Field(i))
		if value == "" {
			continue
		}

		filterParts := strings.Split(filterTag, ",")
		operator := filterParts[0]
		columnName := structField.Name
		for _, part := range filterParts[1:] {
			if strings.HasPrefix(part, "column=") {
				columnName = strings.TrimPrefix(part, "column=")
			}
		}

		filterFields = append(filterFields, FilterField{
			Column:   toColumn(columnName),
			Operator: operator,
			Value:    value,
		})
	}

	return filterFields
}

func FilterDBScopes(db *gorm.DB, fields []FilterField) *gorm.DB {
	for _, field := range fields {
		switch field.Operator {
		case "equal":
			db = db.Where(field.Column+" = ?", field.Value)
		case "like":
			db = db.Where(field.Column+" LIKE ?", "%"+field.Value+"%")
		case "in":
			values := strings.Split(field.Value, ",")
			db = db.Where(field.Column+" IN ?", values)
		case "not_in":
			values := strings.Split(field.Value, ",")
			db = db.Where(field.Column+" NOT IN ?", values)
		case "gt":
			db = db.Where(field.Column+" > ?", field.Value)
		case "lt":
			db = db.Where(field.Column+" < ?", field.Value)
		case "gte":
			db = db.Where(field.Column+" >= ?", field.Value)
		case "lte":
			db = db.Where(field.Column+" <= ?", field.Value)
		case "between":
			parts := strings.Split(field.Value, ",")
			if len(parts) == 2 {
				db = db.Where(field.Column+" BETWEEN ? AND ?", parts[0], parts[1])
			}
		case "not_between":
			parts := strings.Split(field.Value, ",")
			if len(parts) == 2 {
				db = db.Where(field.Column+" NOT BETWEEN ? AND ?", parts[0], parts[1])
			}
		case "is_null":
			db = db.Where(field.Column + " IS NULL")
		case "is_not_null":
			db = db.Where(field.Column + " IS NOT NULL")
		default:
			log.Printf("Unknown query operator: %s", field.Operator)
		}
	}

	return db
}

func toColumn(fieldName string) string {
	var snake []rune
	for i, r := range fieldName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			snake = append(snake, '_', r+('a'-'A'))
		} else {
			if r >= 'A' && r <= 'Z' {
				snake = append(snake, r+('a'-'A'))
			} else {
				snake = append(snake, r)
			}
		}
	}

	return string(snake)
}

func fieldValue(flVal reflect.Value) string {
	if !flVal.IsValid() {
		return ""
	}

	var val string
	switch flVal.Kind() {
	case reflect.String:
		val = flVal.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(flVal.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val = strconv.FormatUint(flVal.Uint(), 10)
	case reflect.Bool:
		if flVal.Bool() {
			return "1"
		}
		return "0"
	case reflect.Slice, reflect.Array:
		var parts []string
		for i := 0; i < flVal.Len(); i++ {
			elem := flVal.Index(i)
			parts = append(parts, fieldValue(elem))
		}
		return strings.Join(parts, ",")
	default:

	}

	return val
}
