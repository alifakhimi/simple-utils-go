package simutils

import (
	"database/sql/driver"
	"encoding/json"
)

// NullBool represents a bool that may be null.
// NullBool implements the [Scanner] interface so
// it can be used as a scan destination, similar to [NullString].
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the [Scanner] interface.
func (n *NullBool) Scan(value any) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true

	if bv, err := driver.Bool.ConvertValue(value); err != nil {
		return err
	} else {
		n.Bool = bv.(bool)
		return nil
	}
}

// Value implements the [driver.Valuer] interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	}
	return json.Marshal(nil)
}

func (n *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}

	if err := json.Unmarshal(b, &n.Bool); err != nil {
		return err
	}

	n.Valid = true
	return nil
}
