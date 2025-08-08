package basemodel

import (
	"encoding/json"
	"fmt"
	"reapp/pkg/context/authctx"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

func (s *SoftFields) BeforeCreate(tx *gorm.DB) error {

	if uid := GetUserID(tx); uid != nil {
		s.CreatedBy = uid
		s.UpdatedBy = uid
	}

	return nil
}

var modelValue map[string]interface{}

func (s *SoftFields) BeforeUpdate(tx *gorm.DB) error {
	if uid := GetUserID(tx); uid != nil {
		s.UpdatedBy = uid
	}

	model := reflect.New(reflect.TypeOf(tx.Statement.Dest).Elem()).Interface()

	pkf := ""
	pkv := interface{}(nil)
	if tx.Statement.Schema != nil {
		for _, field := range tx.Statement.Schema.Fields {
			if field.PrimaryKey {
				pkf = field.DBName
				val := reflect.ValueOf(tx.Statement.Dest)
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				pkv, _ = field.ValueOf(tx.Statement.Context, val)
				break
			}
		}
	}

	if pkf == "" || pkv == nil {
		return fmt.Errorf("could not determine primary key for model")
	}

	if err := tx.Session(&gorm.Session{}).Model(tx.Statement.Dest).Where(fmt.Sprintf("%s = ?", pkf), pkv).First(model).Error; err != nil {
		return fmt.Errorf("fetching original model before update: %w", err)
	}

	if err := SetBfChangeValueFromInstance(tx, model); err != nil {
		return fmt.Errorf("setting model value: %w", err)
	}

	return nil
}

func (s *SoftFields) AfterCreate(tx *gorm.DB) (err error) {
	if tx.Statement.Schema == nil || tx.Statement.Schema.ModelType == nil {
		return nil
	}

	tbName := tx.Statement.Table
	tbID := GetTbID(tx)

	return AddLogData(tx, tbName, tbID, "CREATE", nil, tx.Statement.Dest)
}

func (s *SoftFields) AfterUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Schema == nil || tx.Statement.Schema.ModelType == nil {
		return nil
	}

	beforeChangeData := GetBfChangeValue(tx)
	changes := map[string]map[string]interface{}{}
	for _, field := range tx.Statement.Schema.Fields {
		var current interface{}

		name := field.Name
		if tx.Statement.Dest == nil {
			current = nil
		} else {
			current, _ = field.ValueOf(tx.Statement.Context, reflect.ValueOf(tx.Statement.Dest))
		}
		var oldValue interface{}
		if beforeChangeData != nil {
			oldValue = beforeChangeData[name]
		}

		if !reflect.DeepEqual(current, oldValue) {
			changes[ToColumnName(name)] = map[string]interface{}{
				"from": oldValue,
				"to":   current,
			}
		}
	}

	tbID := GetTbID(tx)

	return AddLogData(tx, tx.Statement.Table, tbID, "UPDATE", changes, tx.Statement.Dest)
}

func GetTbID(tx *gorm.DB) string {
	var tbID string
	if tx.Statement.Schema != nil {
		for _, field := range tx.Statement.Schema.Fields {
			if field.PrimaryKey {
				val := reflect.ValueOf(tx.Statement.Dest)
				if !val.IsValid() {
					continue
				}
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}

				if fieldVal, _ := field.ValueOf(tx.Statement.Context, val); fieldVal != nil {
					var isNonZero bool
					switch v := fieldVal.(type) {
					case int, int8, int16, int32, int64:
						isNonZero = reflect.ValueOf(v).Int() > 0
					case uint, uint8, uint16, uint32, uint64:
						isNonZero = reflect.ValueOf(v).Uint() > 0
					default:
						isNonZero = fieldVal != nil
					}
					if isNonZero {
						tbID = fmt.Sprintf("%v", fieldVal)
					}
				}
				break
			}
		}
	}
	return tbID
}

func AddLogData(tx *gorm.DB, tbName, tbIDStr, action string, changes any, fullModel any) error {
	var changedData []byte
	var fullData []byte
	var err error

	if changes != nil {
		if reflect.TypeOf(changes).Kind() != reflect.Map {
			return fmt.Errorf("changes must be a map, got %s", reflect.TypeOf(changes).Kind())
		}
		changedData, err = json.Marshal(changes)
		if err != nil {
			return err
		}
	}

	fullData, err = json.Marshal(fullModel)
	if err != nil {
		return err
	}
	var createdBy *uint32 = GetUserID(tx)

	var tbID uint64
	if tbIDStr != "" {
		var err error
		tbID, err = strconv.ParseUint(tbIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert modelID to uint64: %w", err)
		}
	}

	log := TableLog{
		TbName:      tbName,
		TbID:        tbID,
		Action:      action,
		ChangedData: changedData,
		FullData:    fullData,
		CreatedBy:   createdBy,
	}

	return tx.Session(&gorm.Session{NewDB: true}).Create(&log).Error
}

func ToColumnName(fieldname string) string {

	var result []rune
	for i, r := range fieldname {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_', r+('a'-'A'))
		} else {
			if r >= 'A' && r <= 'Z' {
				result = append(result, r+('a'-'A'))
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

func GetUserID(tx *gorm.DB) *uint32 {
	return authctx.UserID(tx.Statement.Context)
}
