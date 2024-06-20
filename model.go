package simutils

import (
	"time"

	"gorm.io/gorm"
)

// Model ...
type Model struct {
	ID        PID            `json:"id,omitempty" gorm:"column:id;primarykey"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (m Model) IsDeleted() bool {
	return m.DeletedAt.Valid
}
