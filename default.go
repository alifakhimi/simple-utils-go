package simutils

import "reflect"

// SetToNilIfZeroValue checks if the input value is a zero value. If it is a zero value
// and its type is not a function, channel, or interface, it returns a nil pointer.
// Otherwise, it returns the pointer to the input value.
func SetToNilIfZeroValue[T any](value T) *T {
	// Use reflection to get the value and type of the input
	val := reflect.ValueOf(value)
	kind := val.Kind()

	// Check if the input value is a zero value and its type is not a function, channel, or interface
	if val.IsZero() && kind != reflect.Func && kind != reflect.Chan && kind != reflect.Interface {
		// If the value is zero and not a non-pointer type, return a nil pointer
		return nil
	}

	// If the value is not zero or the type is mutable, return the pointer to the input value
	return &value
}

// DefaultIfZero returns defaultVal if val is the zero value for its type.
// Otherwise, it returns the input value.
func DefaultIfZero[T comparable](val, defaultVal T) T {
	// Use reflection to check if val is the zero value for its type
	if reflect.ValueOf(val).IsZero() {
		// If val is the zero value, return defaultVal
		return defaultVal
	}
	// Otherwise, return the input value
	return val
}
