package simutils

import "errors"

var (
	ErrUserNoFound = errors.New("user/user id not found")
)

type Users []*User

// User ...
type User struct {
	Model
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Username  string   `json:"username" gorm:"unique;not null"`
	Mobile    string   `json:"mobile" gorm:"index"`
	Password  string   `json:"-"`
	Email     string   `json:"email"`
	Status    UserMode `json:"status" gorm:"default:1"`
}

func (u *User) FullName() string {
	return u.Firstname + " " + u.Lastname
}
