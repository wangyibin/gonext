package apidoc

import (
	"fmt"
	"reflect"
	"strconv"
)

func toTargeType(targetType reflect.Type, value string) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint(i)), err
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Ptr:
		if len(value) == 0 {
			return reflect.Zero(targetType), nil
		}
		switch targetType.Elem().Kind() {
		case reflect.Int:
			i, err := strconv.Atoi(value)
			return reflect.ValueOf(&i), err
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			i, err := strconv.ParseUint(value, 10, 64)
			ui := uint(i)
			return reflect.ValueOf(&ui), err
		case reflect.String:
			return reflect.ValueOf(&value), nil
		}
	}
	return reflect.Zero(targetType), fmt.Errorf("unsupport param type: %s", targetType.Kind())
}
