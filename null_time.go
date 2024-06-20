package simutils

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NullTime sql.NullTime

// Scan implements the Scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}

	if err := json.Unmarshal(b, &n.Time); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
