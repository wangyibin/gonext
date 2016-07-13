package typeconv

import (
	"fmt"
	"reflect"
	"strconv"
)

func SetValue(field reflect.Value, value string) error {
	targetType := field.Type()
	switch targetType.Kind() {
	case reflect.Ptr:
		return SetValue(field.Elem(), value)
		//v, err := ToTargetType(targetType.Elem(), value)
		//return v.Pointer(), err
	default:
		v, err := extractBaseTypeValue(targetType, value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(v))
		return nil
		//return reflect.ValueOf(v), err
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