package typeconv

import (
	"fmt"
	"reflect"
	"strconv"
)

func ToTargetType(targetType reflect.Type, value string) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.Ptr:
		v, err := ToTargetType(targetType.Elem(), value)
		return reflect.ValueOf(&v), err
	default:
		v, err := extractBaseTypeValue(targetType, value)
		return reflect.ValueOf(v), err
	}
}

func extractBaseTypeValue(targetType reflect.Type, value string) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		return int(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		return int8(i), err
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		return int16(i), err
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		return int32(i), err
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		return i, err
	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 64)
		return uint(i), err
	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 8)
		return uint8(i), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 16)
		return uint16(i), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 32)
		return uint32(i), err
	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		return i, err
	case reflect.Uintptr:
		i, err := strconv.ParseUint(value, 10, 64)
		return uint(i), err
	case reflect.String:
		return value, nil
	}
	return reflect.Zero(targetType), fmt.Errorf("unsupport param type: %s", targetType.Kind())
}