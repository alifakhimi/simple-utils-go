package simutils

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrRecordNotFound   = errors.New("record not found")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrAlreadyExist     = errors.New("already exist")
	ErrSystemItemDelete = errors.New("you cant delete system items")
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}
