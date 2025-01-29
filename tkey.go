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

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

var (
	tkeyRegex = regexp.MustCompile(`^(\w+):(\w+(?:,\w+)*)$`) // Matches keys in the format <table_name>:<value1,value2,...>.
)

type TKey string

// IsValid checks if the TKey matches the required format.
func (k TKey) IsValid() bool {
	return k != "" && tkeyRegex.MatchString(string(k))
}

// GetTKey generates a unique key for a struct based on its primary key fields and table name.
func GetTKey(val any, depth ...int) TKey {
	if len(depth) == 0 {
		depth = append(depth, 0)
	}

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
			fieldValue := v.Field(i)

			// If the field is an embedded struct, process it recursively.
			if f.Anonymous {
				embeddedKey := GetTKey(fieldValue.Interface(), depth[0]+1)
				if embeddedKey != "" {
					keys = append(keys, string(embeddedKey))
				}
				continue
			}

			// Check for `sim` tag with primaryKey.
			if AnyItemExists(strings.Split(strings.ToLower(f.Tag.Get("sim")), ";"), []string{"primarykey", "primary_key"}) {
				if val := fmt.Sprintf("%v", v.Field(i)); val != "" {
					simKeys = append(simKeys, val)
				}
			} else if AnyItemExists(strings.Split(strings.ToLower(f.Tag.Get("gorm")), ";"), []string{"primarykey", "primary_key"}) {
				if val := fmt.Sprintf("%v", v.Field(i)); val != "" {
					gormKeys = append(gormKeys, val)
				}
			}
		}

		if len(simKeys) > 0 {
			keys = append(keys, simKeys...)
		} else if len(gormKeys) > 0 {
			keys = append(keys, gormKeys...)
		}

		if depth[0] == 0 {
			return TKey(fmt.Sprintf("%s:%s", GetModelName(val), strings.Join(keys, ",")))
		} else {
			return TKey(strings.Join(keys, ","))
		}
	}

	return TKey("")
}

// GetTableName retrieves the table name for a given struct using GORM v1 conventions.
func GetModelName(val interface{}) string {
	v := reflect.TypeOf(val)

	// If it's a pointer, get the underlying type
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Ensure it's a struct
	if v.Kind() != reflect.Struct {
		return ""
	}

	// Default behavior: Convert struct name to snake_case and pluralize
	return inflection.Plural(strcase.ToSnake(v.Name()))
}
