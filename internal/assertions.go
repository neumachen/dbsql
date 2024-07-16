package internal

import "reflect"

func IsNil(v any) bool {
	if v == nil {
		return true
	}

	valueOf := reflect.ValueOf(v)

	switch valueOf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return valueOf.IsNil()
	}

	return false
}

func IsNilOrZeroValue(v any) bool {
	if v == nil {
		return true
	}

	valueOf := reflect.ValueOf(v)

	switch valueOf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return valueOf.IsNil() || valueOf.IsZero()
	case reflect.Int:
		return valueOf.IsZero()
	default:
		return false
	}
}

// IsNil checks if a value is nil.
// func IsNil[T any](v T) bool {
// 	// Use type switch to handle different types
// 	switch val := any(v).(type) {
// 	case nil:
// 		return true
// 	case interface{ IsNil() bool }:
// 		return val.IsNil()
// 	default:
// 		return false
// 	}
// }
//
// // IsNilOrZeroValue checks if a value is nil or a zero value.
// func IsNilOrZeroValue[T comparable](v T) bool {
// 	// First, check if it's nil
// 	if IsNil(v) {
// 		return true
// 	}
//
// 	// Then, check if it's a zero value
// 	var zero T
// 	return v == zero
// }
//
// // IsNilOrZeroValueSlice checks if a slice is nil or has zero length.
// // NOTE: For types that don't support direct comparison (like slices)
// // you might need a separate function.
// func IsZeroValueSlice[T any](v []T) bool {
// 	return len(v) == 0
// }
//
// // IsZeroValueMap checks if a map has zero length.
// func IsZeroValueMap[K comparable, V any](v map[K]V) bool {
// 	return len(v) == 0
// }
//
// func IsNilChan[T any](v chan T) bool {
// 	return v == nil
// }
