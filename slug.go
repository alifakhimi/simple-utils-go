package simutils

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cast"
)

// Errors
var (
	ErrInvalidSlug = errors.New("invalid slug")
)

var (
	// Regular expression for a slug that supports Persian, English, numbers, and dashes
	slugRegex = regexp.MustCompile(`^[\p{L}\p{N}-]+$`)
)

type Slug string

func Invalidate(s *Slug) *Slug {
	if s == nil {
		return nil
	}

	return MakeSlugP(ToString(s))
}

// IsSlug checks if a string is a valid slug
func IsSlug(a any) bool {
	if a == nil {
		return false
	}

	// Use reflection to check if the value is nil
	val := reflect.ValueOf(a)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return false
	}

	// Check if the string matches the slug pattern
	return slugRegex.MatchString(cast.ToString(a))
}

// DecodeSlug decodes a URL-encoded string and checks if it's a valid slug
func DecodeAndCheckSlug(s string) (string, error) {
	// Decode the URL-encoded string
	decodedStr, err := url.QueryUnescape(s)
	if err != nil {
		return "", err
	}

	// Convert to lowercase
	decodedStr = strings.ToLower(decodedStr)

	// Check if the decoded string is a valid slug
	if !IsSlug(decodedStr) {
		return "", ErrInvalidSlug
	}

	return decodedStr, nil
}

// MakeSlug takes any string and converts it to a URL-friendly slug
func MakeSlug(name string) Slug {
	s, _ := MakeSlugE(name)
	return s
}

// MakeSlugP takes any string and converts it to a URL-friendly *slug
func MakeSlugP(name string) *Slug {
	if s, err := MakeSlugE(name); err != nil {
		return nil
	} else {
		return &s
	}
}

// MakeSlugE takes any string and converts it to a URL-friendly slug
func MakeSlugE(name string) (Slug, error) {
	if IsURLEncoded(name) {
		if s, err := DecodeAndCheckSlug(name); err != nil {
			return "", err
		} else {
			name = s
		}
	}

	// Trim spaces from start and end of the string
	name = strings.TrimSpace(name)

	// Replace spaces with dashes
	name = strings.ReplaceAll(name, " ", "-")

	// Remove any character that is not a letter, number, or dash
	re := regexp.MustCompile(`[^\p{L}\p{N}-]`)
	name = re.ReplaceAllString(name, "")

	// Convert to lower case
	name = strings.ToLower(name)

	// Remove multiple dashes
	name = strings.ReplaceAll(name, "--", "-")

	if len(name) == 0 || name == "-" {
		return "", ErrInvalidSlug
	}

	return Slug(name), nil
}

func ToString(a any) string {
	switch s := a.(type) {
	case Slug:
		return string(s)
	case *Slug:
		return string(*s)
	default:
		return ""
	}
}

func (s Slug) IsValid() bool {
	return IsSlug(string(s))
}

func (s Slug) Append(a any) Slug {
	return MakeSlug(fmt.Sprintf("%s-%v", s, a))
}
