package mapper

import (
	"errors"
	"reflect"
)

func CloneStructFields(sourceStruct interface{}, fieldNames []string) map[string]interface{} {
	sourceValue := reflect.Indirect(reflect.ValueOf(sourceStruct))
	sourceType := sourceValue.Type()

	fields := make(map[string]interface{})
	for idx := 0; idx < sourceType.NumField(); idx++ {
		field := sourceType.Field(idx)
		if field.PkgPath != "" {
			continue
		}
		for _, fieldName := range fieldNames {
			if field.Name == fieldName {
				fields[field.Name] = sourceValue.Field(idx).Interface()
				break
			}
		}
	}

	return fields
}

func MapStructFields(targetStruct interface{}, fields map[string]interface{}) error {
	targetValue := reflect.ValueOf(targetStruct)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	for fieldName, fieldValue := range fields {
		field := targetValue.FieldByName(fieldName)
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

// MapStruct copies selected fields from sourceStruct to targetStruct.
//
// Parameters:
// - targetStruct: a pointer to the struct that will receive the field values.
// - sourceStruct: the struct to copy values from.
// - fieldNames: a slice of strings specifying which field names to copy.
//
// Returns:
// - error: if targetStruct is not a pointer, or if type conversion fails.
//
// Notes:
// - Fields must be exported (start with an uppercase letter) to be copied.
// - Only fields listed in fieldNames and present in both structs will be copied.
func MapStruct(targetStruct interface{}, sourceStruct interface{}, fieldNames []string) error {
	targetValue := reflect.ValueOf(targetStruct)
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("the destination data must be provided as a reference")
	}
	fields := CloneStructFields(sourceStruct, fieldNames)
	return MapStructFields(targetStruct, fields)
}
