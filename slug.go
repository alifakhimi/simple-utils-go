package simutils

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var (
	ErrInvalidSlug = errors.New("invalid slug")
)

type Slug string

func (s Slug) IsValid() bool {
	return IsSlug(string(s))
}

// IsSlug checks if a string is a valid slug
func IsSlug(s string) bool {
	// Regular expression for a slug that supports Persian, English, numbers, and dashes
	slugRegex := regexp.MustCompile(`^[\p{L}\p{N}-]+$`)

	// Check if the string matches the slug pattern
	return slugRegex.MatchString(s)
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
func MakeSlug(name string) (Slug, error) {
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

	if len(name) < 3 {
		return "", ErrInvalidSlug
	}

	return Slug(name), nil
}

func (s Slug) ToString() string {
	return string(s)
}
