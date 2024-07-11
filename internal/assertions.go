package internal

import (
	"reflect"
)

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
