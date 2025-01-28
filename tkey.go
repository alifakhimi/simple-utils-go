/*
Package simutils provides utilities for generating unique textual keys (TKey) from structs,
using specific tags (e.g., `sim` or `gorm`) to identify primary key fields.

Key Features:
1. **TKey Type**:
  - A string wrapper representing a unique key in the format `<table_name>:<primary_keys>`.

2. **Regex Validation**:
  - Ensures keys follow a strict pattern using regular expressions.

3. **GetTKey Function**:
  - Extracts primary key fields based on `sim` or `gorm` tags using reflection.
  - Combines table name and primary key values into a TKey.

Usage:
- Use `GetTKey` to generate unique keys for caching or database operations.
- Use `IsValid` to validate the format of a key.
*/

package simutils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	tkeyRegex = regexp.MustCompile(`^(\w+):(\w+(?:,\w+)*)$`) // Matches keys in the format <table_name>:<value1,value2,...>.
)

type TKey string

// IsValid checks if the TKey matches the required format.
func (k TKey) IsValid() bool {
	return tkeyRegex.MatchString(string(k))
}

// GetTKey generates a unique key for a struct based on its primary key fields and table name.
func GetTKey(val any) TKey {
	keys := []string{}
	v := reflect.ValueOf(val)

	// Handle pointers by dereferencing them.
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if v.Kind() == reflect.Struct {
		var (
			gormKeys = []string{}
			simKeys  = []string{}
		)

		for i := 0; i < v.NumField(); i++ {
			t := v.Type()
			f := t.Field(i)

			// Check for `sim` tag with primaryKey.
			if ItemExists(strings.Split(strings.ToLower(f.Tag.Get("sim")), ";"), strings.ToLower("primaryKey")) {
				simKeys = append(simKeys, fmt.Sprintf("%v", v.Field(i)))
			} else if ItemExists(strings.Split(strings.ToLower(f.Tag.Get("gorm")), ";"), strings.ToLower("primaryKey")) {
				gormKeys = append(gormKeys, fmt.Sprintf("%v", v.Field(i)))
			}
		}

		if len(simKeys) > 0 {
			keys = append(keys, simKeys...)
		} else if len(gormKeys) > 0 {
			keys = append(keys, gormKeys...)
		}

		return TKey(fmt.Sprintf("%s:%s", GetTableName(val), strings.Join(keys, ",")))
	}

	return TKey("")
}
