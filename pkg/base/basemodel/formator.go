package basemodel

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

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

func (d *TDateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = TDateOnly{Time: time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = TDateOnly{Time: v}
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*d = TDateOnly{Time: t}
	default:
		return fmt.Errorf("cannot scan type %T into TDateOnly", value)
	}
	return nil
}

func (d TDateOnly) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Format("2006-01-02"), nil
}

func (d TDateOnly) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format("2006-01-02")
}

func (d TDateOnly) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", d.Format("2006-01-02"))
	return []byte(formatted), nil
}

func (d *TDateOnly) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" || str == "null" {
		*d = TDateOnly{Time: time.Time{}}
		return nil
	}
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", str)
	}
	d.Time = t
	return nil
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

// Write value to DB
func (t TString) Value() (driver.Value, error) {
	s := strings.TrimSpace(string(t))
	if s == "" {
		return nil, nil
	}
	return s, nil
}

// Read value from DB
func (t *TString) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	if b, ok := value.([]byte); ok {
		*t = TString(string(b))
		return nil
	}
	*t = TString(value.(string))
	return nil
}
