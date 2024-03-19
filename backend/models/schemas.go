package models

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

// "2023-01-22 02:03:51"
// will be converted to
// 2023-01-22 10:03:51.000

type MyTime struct {
	time.Time
}

var formatTime = "2006-01-02 15:04:05"

func (t *MyTime) UnmarshalJSON(b []byte) (err error) {
	b = bytes.Trim(b, "\"")
	tt, err := time.Parse(formatTime, string(b))
	*t = MyTime{tt}
	//now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(b), time.Local)
	return
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Format(formatTime))), nil
}

func (ct MyTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

func (ct *MyTime) Scan(value interface{}) error {
	if value == nil {
		*ct = MyTime{time.Time{}}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("could not convert value to time.Time")
	}
	*ct = MyTime{t}
	return nil
}
