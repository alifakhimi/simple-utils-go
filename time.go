package simutils

import (
	"database/sql/driver"
	"time"
)

// Time ...
type Time struct {
	time.Time
}

// NowTime returns the current local time.
func NowTime() Time {
	return Time{Time: time.Now()}
}

// Value - Implementation of valuer for database/sql
func (t Time) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	// such as string, bool and ...
	return t.Format(time.RFC3339), nil
}

// Scan ...
func (t *Time) Scan(value interface{}) error {
	var err error
	switch v := value.(type) {
	case string:
		var tm time.Time
		tm, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		*t = Time{tm}
	case Time:
		*t = v
	case time.Time:
		*t = Time{v}
	}
	return err
}

// TODO complete here to use nulltime work
