package simutils

import (
	"errors"
	"time"

	"gorm.io/gorm"
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
	ID          PID            `json:"id,omitempty" gorm:"column:id;primaryKey;"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;"`
	Active      NullBool       `json:"active,omitempty" sql:"default:true;"`
	Version     int            `json:"version,omitempty" sql:"default:NULL;"`
	Description string         `json:"description,omitempty"`
	UserID      PID            `json:"user_id,omitempty" sql:"default:NULL;"`
	User        *User          `json:"user,omitempty" gorm:"<-:false;"`
	Meta        JSON           `json:"meta,omitempty" sql:"default:NULL;"`
	Error       Error          `json:"error,omitempty" gorm:"-"`
}

func (m CommonTableFields) IsDeleted() bool {
	return m.DeletedAt.Valid
}

type PolymorphicFields struct {
	OwnerType string `json:"owner_type,omitempty"`
	OwnerID   PID    `json:"owner_id,omitempty" gorm:"index"`
}

// ToCommonTableFields to common table fields
func ToCommonTableFields(id PID, user *User) (ctf CommonTableFields, err error) {
	ctf = CommonTableFields{
		ID:   id,
		User: user,
	}

	if user == nil {
		err = ErrCommonTableFieldsInvalidUser
	} else {
		ctf.UserID = user.ID
	}

	return
}

func CheckNilError(values ...interface{}) error {
	for _, value := range values {
		if value == nil {
			return errors.New("unexpected nil value")
		}
	}
	return nil
}
