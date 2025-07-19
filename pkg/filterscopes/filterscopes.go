package filterscopes

import (
	"log"
	"net/url"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type QueryFilter struct {
	Field    string
	Operator string
	Value    string
}

func ParseQueryByUrlValues(query interface{}, values url.Values) []QueryFilter {
	var filters []QueryFilter
	valOfQuery := reflect.ValueOf(query)
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
		formFieldKey := structField.Tag.Get("form")
		if formFieldKey == "" {
			formFieldKey = structField.Name
		}
		if value := values.Get(formFieldKey); value != "" {
			filterParts := strings.Split(filterTag, ",")
			filterOperator := filterParts[0]
			columnName := structField.Name
			for _, part := range filterParts[1:] {
				if strings.HasPrefix(part, "column=") {
					columnName = strings.TrimPrefix(part, "column=")
				}
			}
			filters = append(filters, QueryFilter{
				Field:    columnName,
				Operator: filterOperator,
				Value:    value,
			})
		}
	}
	return filters
}

func QueryFilterScopes(db *gorm.DB, qfs []QueryFilter) *gorm.DB {
	for _, filter := range qfs {
		switch filter.Operator {
		case "equal":
			db = db.Where(ToColumn(filter.Field)+" = ?", filter.Value)
		case "like":
			db = db.Where(ToColumn(filter.Field)+" LIKE ?", "%"+filter.Value+"%")
		case "in":
			values := strings.Split(filter.Value, ",")
			db = db.Where(ToColumn(filter.Field)+" IN ?", values)
		case "not_in":
			values := strings.Split(filter.Value, ",")
			db = db.Where(ToColumn(filter.Field)+" NOT IN ?", values)
		case "gt":
			db = db.Where(ToColumn(filter.Field)+" > ?", filter.Value)
		case "lt":
			db = db.Where(ToColumn(filter.Field)+" < ?", filter.Value)
		case "gte":
			db = db.Where(ToColumn(filter.Field)+" >= ?", filter.Value)
		case "lte":
			db = db.Where(ToColumn(filter.Field)+" <= ?", filter.Value)
		case "between":
			parts := strings.Split(filter.Value, ",")
			if len(parts) == 2 {
				db = db.Where(ToColumn(filter.Field)+" BETWEEN ? AND ?", parts[0], parts[1])
			}
		case "not_between":
			parts := strings.Split(filter.Value, ",")
			if len(parts) == 2 {
				db = db.Where(ToColumn(filter.Field)+" NOT BETWEEN ? AND ?", parts[0], parts[1])
			}
		case "is_null":
			db = db.Where(ToColumn(filter.Field) + " IS NULL")
		case "is_not_null":
			db = db.Where(ToColumn(filter.Field) + " IS NOT NULL")
		default:
			log.Printf("Unknown query operator: %s", filter.Operator)
		}
	}
	return db
}

func ToColumn(fieldName string) string {
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
