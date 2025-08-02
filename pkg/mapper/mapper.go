package mapper

import (
	"reflect"
	"strings"
)

func GetFieldsOfDTO(modelDTO any, fieldNames []string) map[string]any {
	dtoValue := reflect.Indirect(reflect.ValueOf(modelDTO))
	dtoType := dtoValue.Type()

	fields := make(map[string]any)
	for idx := 0; idx < dtoType.NumField(); idx++ {
		field := dtoType.Field(idx)
		if field.PkgPath != "" {
			continue
		}
		for _, fieldName := range fieldNames {
			if strings.EqualFold(field.Name, fieldName) {
				fields[field.Name] = dtoValue.Field(idx).Interface()
				break
			}
		}
	}

	return fields
}

func AssignFieldValuesToModel(model any, fields map[string]any) error {
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	for fieldName, fieldValue := range fields {
		field := modelValue.FieldByName(fieldName)
		if !field.IsValid() {
			return nil
		}
		if !field.CanSet() {
			continue
		}

		value := reflect.ValueOf(fieldValue)
		if value.Type().ConvertibleTo(field.Type()) {
			field.Set(value.Convert(field.Type()))
		} else {
			return nil
		}
	}

	return nil
}

func AssignModelValues[T any](model *T, modelDTO any, fieldNames []string) error {
	fields := GetFieldsOfDTO(modelDTO, fieldNames)
	return AssignFieldValuesToModel(model, fields)
}
