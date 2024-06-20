package simutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidPID = errors.New("id is not valid")
)

type PIDs []PID

// PID Primary ID
type PID int64

// ToPID create pid
func ToPID(id int) PID {
	if id < 0 {
		id = 0
	}

	return PID(id)
}

// NilPID Null Primary ID
var NilPID = PID(0)

// NullPID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
type NullPID struct {
	PID   PID
	Valid bool
}

// ToNullPID new instance of NullPID with checking structure
func (id *PID) ToNullPID() (NullPID, error) {
	if !IsValid(id) {
		return NullPID{}, ErrInvalidPID
	}

	return NullPID{
		Valid: true,
		PID:   *id,
	}, nil
}

// Value - Implementation of valuer for database/sql
func (id PID) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	// such as string, bool and ...
	return int64(id), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice is handled by UnmarshalBinary, while
// a longer byte slice or a string is handled by UnmarshalText.
func (id *PID) Scan(src interface{}) error {
	if src == nil {
		*id = PID(0)
		return nil
	}

	// ns := sql.NullInt64{}
	// if err := ns.Scan(src); err != nil {
	//     return err
	// }
	//
	// if !ns.Valid {
	//     return errors.New("scan not valid")
	// }
	//
	// nsv, _ := ns.Value()
	// *id = PID(nsv.(int64))

	*id = PID(src.(int64))

	return nil
}

func (id PID) String() string {
	return strconv.Itoa(int(id))
}

// CheckPID ...
func (id PID) CheckPID() bool {
	return true
}

func (id PID) IsValid() bool {
	return int64(id) > 0
}

func (id NullPID) IsValid() bool {
	if id.Valid {
		return id.PID.IsValid()
	}
	return false
}

// TODO implement NullPID type
// IsValid check id validation
func IsValid(id interface{}) bool {
	if id == nil || id == (*PID)(nil) || !reflect.ValueOf(id).IsValid() {
		return false
	}

	val := reflect.ValueOf(id)
	if val.Kind() == reflect.Ptr {
		return (val.Interface().(*PID)).IsValid()
	}

	return Parse(id).IsValid()
}

// ParsePID , parses a string id to a PID one
func ParsePID(id interface{}) (pid PID, err error) {
	switch id := id.(type) {
	case string:
		d, _ := strconv.ParseInt(id, 10, 64)
		if strings.EqualFold(strconv.FormatInt(d, 10), id) {
			pid = PID(d)
		}
	case int:
		pid = PID(id)
	case int64:
		pid = PID(id)
	case float64:
		pid = PID(id)
	case PID:
		pid = id
	}

	if !pid.IsValid() {
		err = ErrInvalidPID
	}

	return pid, err
}

// Parse ...
func Parse(id interface{}) PID {
	pid, _ := ParsePID(id)
	return pid
}

// Validate ...
func Validate(id string) (PID, bool) {
	pid, err := ParsePID(id)
	return pid, err == nil
}

// String ...
func String(id PID) string {
	return id.String()
}

// CheckPID ...
func CheckPID(id PID) bool {
	return id.CheckPID()
}

// Value implements the driver.Valuer interface.
func (u NullPID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	// Delegate to int64 Value function
	return u.PID.Value()
}

// Scan implements the sql.Scanner interface.
func (u *NullPID) Scan(src interface{}) error {
	if src == nil {
		u.PID, u.Valid = NilPID, false
		return nil
	}

	// Delegate to int64 Scan function
	u.Valid = true
	return u.PID.Scan(src)
}

// MarshalJSON ...
func (u NullPID) MarshalJSON() ([]byte, error) {
	if u.Valid {
		return json.Marshal(u.PID)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON ...
func (u *NullPID) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *PID
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		u.Valid = true
		u.PID = *x
	} else {
		u.Valid = false
	}
	return nil
}

// GetPIDsDefault ...
func GetPIDsDefault(s []string, def []PID) (pids []PID) {
	if len(s) > 0 {
		for _, sid := range s {
			if id, err := ParsePID(sid); err == nil {
				pids = append(pids, id)
			}
		}

		return pids
	}

	return def
}

func GetStringFromPIDs(pids []PID) (strPIDs []string) {
	stringPIDs := make([]string, len(pids))
	for i, pid := range pids {
		stringPIDs[i] = pid.String()
	}
	return stringPIDs
}

func AppendIfNotExist(s []PID, e PID) []PID {
	for _, a := range s {
		if a == e {
			return s
		}
	}
	return append(s, e)
}

func IDExist(s []PID, e PID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IDsToString(ids []PID) (s string) {
	if len(ids) == 0 {
		return s
	}

	for _, id := range ids {
		s += id.String() + ","
	}

	return s[:len(s)-1]
}

func StringToIDs(s string) (ids []PID) {
	for _, id := range strings.Split(s, ",") {
		ids = append(ids, Parse(id))
	}

	return
}
