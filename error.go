package simutils

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrRecordNotFound   = errors.New("record not found")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrAlreadyExist     = errors.New("already exist")
	ErrSystemItemDelete = errors.New("you cant delete system items")
)
