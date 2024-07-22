package simutils

import (
	"database/sql"
	"errors"
)

var (
	// ErrCommonTableFieldsInvalidID invalid id
	ErrCommonTableFieldsInvalidID = errors.New("invalid id")
	// ErrCommonTableFieldsInvalidUser invalid user
	ErrCommonTableFieldsInvalidUser = errors.New("invalid user")
	// ErrCommonTableFieldsIDNotMatch id is not match
	ErrCommonTableFieldsIDNotMatch = errors.New("id is not match")
)

// CommonTableFields common table fields
type CommonTableFields struct {
	Model
	Active      bool           `json:"active,omitempty" gorm:"default:false"`
	Version     int64          `json:"version,omitempty"`
	Description sql.NullString `json:"description,omitempty"`
	UserID      NullPID        `json:"user_id,omitempty"`
	User        *User          `json:"user,omitempty" gorm:"<-:false"`
	Meta        JSON           `json:"meta,omitempty"`
	Error       Error          `json:"error,omitempty" gorm:"-"`
}

type PolymorphicFields struct {
	OwnerType string `json:"owner_type,omitempty"`
	OwnerID   PID    `json:"owner_id,omitempty" gorm:"index"`
}

// ToCommonTableFields to common table fields
func ToCommonTableFields(id PID, user *User) (ctf CommonTableFields, err error) {
	if user == nil {
		return ctf, ErrCommonTableFieldsInvalidUser
	}

	userNullPID, err := user.ID.ToNullPID()
	if err != nil {
		return ctf, err
	}

	return CommonTableFields{
		Model: Model{
			ID: id,
		},
		UserID: userNullPID,
		User:   user,
	}, nil
}

func CheckNilError(values ...interface{}) error {
	for _, value := range values {
		if value == nil {
			return errors.New("unexpected nil value")
		}
	}
	return nil
}
