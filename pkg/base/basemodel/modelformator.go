package basemodel

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type TDateTime struct {
	time.Time
}

type TString string

func (dtf *TDateTime) Scan(value interface{}) error {
	if value == nil {
		*dtf = TDateTime{Time: time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*dtf = TDateTime{Time: v}
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*dtf = TDateTime{Time: t}
	default:
		return fmt.Errorf("cannot scan type %T into DateTimeFormat", value)
	}
	return nil
}

func (dtf TDateTime) Value() (driver.Value, error) {
	return dtf.Time, nil
}

func (dtf TDateTime) String() string {
	return dtf.Format("2006-01-02 15:04:05")
}

func (dtf TDateTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", dtf.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (dtf *TDateTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = strings.Trim(str, `"`)
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, str)
		if err == nil {
			dtf.Time = t
			return nil
		}
	}
	return fmt.Errorf("invalid time format: %s", str)
}

func (t *TString) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	*t = TString(strings.TrimSpace(s))
	return nil
}

func (t *TString) UnmarshalParam(param string) error {
	*t = TString(strings.TrimSpace(param))
	return nil
}
