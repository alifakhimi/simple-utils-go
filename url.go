package simutils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

type URL struct {
	url.URL
}

func (u *URL) String() string {
	return u.URL.String()
}

func (u *URL) IsValid() bool {
	return govalidator.IsURL(u.String())
}

func (u *URL) Clone() *URL {
	return cloneURL(u)
}

// TODO refactor with regexp
// LastRoute returns last route of url path
//
// LastRoute(false): https://sample.com/path/to/any/ => "")
//
// LastRoute(true): https://sample.com/path/to/any/ => "any")
func (u *URL) LastRoute(removeTrailingSlash bool) string {
	p := u.Path
	if removeTrailingSlash {
		p = strings.TrimSuffix(p, "/")
	}
	paths := strings.Split(p, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[len(paths)-1]
}

// MarshalJSON to output non base64 encoded []byte
func (u URL) MarshalJSON() ([]byte, error) {
	if !u.IsValid() {
		return json.Marshal(nil)
	}
	return json.Marshal(u.String())
}

// UnmarshalJSON to deserialize []byte
func (u *URL) UnmarshalJSON(b []byte) (err error) {
	var (
		rawURL    string
		parsedURL *url.URL
	)

	if err = json.Unmarshal(b, &rawURL); err != nil {
		return err
	} else if rawURL == "" {
		return nil
	}

	if parsedURL, err = url.ParseRequestURI(rawURL); err != nil {
		return err
	}

	*u = *FromURL(parsedURL)

	return nil
}

// GormDataType gorm common data type
func (URL) GormDataType() string {
	return "string"
}

func (u *URL) Scan(b interface{}) (err error) {
	if b == nil {
		return nil
	}
	text, ok := b.(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal 'string' value: %v", b)
	}
	if text == "" {
		return nil
	}
	v, err := url.Parse(text)
	if err != nil {
		return err
	}
	*u = *FromURL(v)
	return nil
}

func (u URL) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	// such as string, bool and ...
	if !u.IsValid() {
		return nil, nil
	}
	return u.String(), nil
}

func URLFromString(rawURL string) *URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	return FromURL(u)
}

func FromURL(u *url.URL) *URL {
	if u == nil {
		return nil
	}

	return &URL{
		URL: *u,
	}
}

func cloneURL(u *URL) *URL {
	if u == nil {
		return nil
	}
	u2 := new(URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

// IsURLEncoded checks if a string contains URL-encoded characters
func IsURLEncoded(s string) bool {
	// Regular expression to detect URL-encoded characters (e.g., %20, %D8, etc.)
	encodedRegex := regexp.MustCompile(`%[0-9A-Fa-f]{2}`)

	// Check if the string contains URL-encoded characters
	return encodedRegex.MatchString(s)
}
