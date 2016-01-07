package apidoc

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/labstack/echo"
)

// BuildEchoHandler func
func BuildEchoHandler(fullRequestPath string, h1 interface{}, h2 interface{}, h3 interface{}) echo.HandlerFunc {
	inType, _, _ := validateChain(h1, h2, h3)

	return func(c *echo.Context) error {
		var requestObj reflect.Value

		if inType != nil {
			requestObj = newType(fullRequestPath, inType, c)
		}
		inParams := make(map[reflect.Type]reflect.Value)
		inParams[inType] = requestObj
		inParams[reflect.TypeOf(c)] = reflect.ValueOf(c)

		lastHandler := h1
		out, err := callHanlder(h1, inParams)
		if err != nil {
			return err
		}
		if h2 != nil {
			lastHandler = h2
			out, err = callHanlder(h2, inParams)
			if err != nil {
				return err
			}
		}

		if h3 != nil {
			lastHandler = h3
			out, err = callHanlder(h3, inParams)
			if err != nil {
				return err
			}
		}
		if len(out) > 1 {
			return fmt.Errorf("return more then one data value is not supported: %s", runtime.FuncForPC(reflect.ValueOf(lastHandler).Pointer()).Name())
		} else if len(out) == 0 {
			return c.NoContent(http.StatusOK)
		} else {
			return c.JSON(http.StatusOK, out[0].Interface())
		}
	}
}

func callHanlder(handler interface{}, inParams map[reflect.Type]reflect.Value) ([]reflect.Value, error) {
	handlerRef := reflect.ValueOf(handler)
	var params []reflect.Value
	for i := 0; i < handlerRef.Type().NumIn(); i++ {
		params = append(params, inParams[handlerRef.Type().In(i)])
	}
	values := handlerRef.Call(params)

	var notErrors []reflect.Value
	var err error
	for _, value := range values {
		if value.Interface() != nil {
			if !isErrorType(value) {
				inParams[value.Type()] = value
				notErrors = append(notErrors, value)
			} else {
				err = value.Interface().(error)
			}
		}
	}
	return notErrors, err
}

func isErrorType(v reflect.Value) bool {
	return v.MethodByName("Error").IsValid()
}
func newType(fullRequestPath string, typ reflect.Type, c *echo.Context) reflect.Value {
	requestType := typ
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	requestObj := reflect.New(requestType)

	pnames := PathNames(fullRequestPath)

	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)
		isPathParam := false
		for _, pname := range pnames {
			if pname == lowCamelStr(field.Name) {
				isPathParam = true
				break
			}
		}
		if field.Name != "Body" {
			var value string
			if isPathParam {
				value = c.Param(lowCamelStr(field.Name))
			} else {
				value = c.Query(lowCamelStr(field.Name))
			}
			setValue(requestObj.Elem().FieldByName(field.Name), field.Name, value)
		} else {
			bodyType := field.Type
			var body interface{}
			if bodyType.Kind() == reflect.Ptr {
				body = reflect.New(field.Type.Elem()).Interface()
			} else {
				body = reflect.New(field.Type).Interface()
			}

			c.Bind(body)
			fmt.Printf("%s\n", body)
			if bodyType.Kind() == reflect.Ptr {
				requestObj.Elem().FieldByName("Body").Set(reflect.ValueOf(body))
			} else {
				requestObj.Elem().FieldByName("Body").Set(reflect.ValueOf(body).Elem())
			}
		}

	}
	return requestObj
}

func setValue(field reflect.Value, name string, value string) error {
	v, err := toTargeType(field.Type(), value)
	fmt.Printf("setValue [%s] -> %s(%v)\n", name, v, v.Type())
	field.Set(v)

	return err
}
