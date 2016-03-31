package gonext

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v8"
	"errors"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

// BuildEchoHandler func
func BuildEchoHandler(fullRequestPath string, handlers []interface{}) echo.HandlerFunc {
	inTypes, _, _ := validateChain(handlers)

	return func(echoContext echo.Context) error {
		// var requestObj reflect.Value
		var err error
		var c Context
		c = &echoContextWrapper{c: echoContext}
		inParams := make(map[reflect.Type]reflect.Value)
		inParams[reflect.TypeOf(c)] = reflect.ValueOf(c)
		for _, inType := range inTypes {
			requestObj, err := newType(fullRequestPath, inType, c)
			if err != nil {
				return err
			}
			inParams[inType] = requestObj
		}

		var lastHandler interface{}
		var out []reflect.Value

		fmt.Printf("call %s\n", fullRequestPath)
		for inParamKey := range inParams {
			fmt.Printf("    in[%s]\n", inParamKey)
		}
		for _, h := range handlers {
			lastHandler = h
			out, err = callHandler(h, inParams)
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

func callHandler(handler interface{}, inParams map[reflect.Type]reflect.Value) ([]reflect.Value, error) {
	handlerRef := reflect.ValueOf(handler)
	var params []reflect.Value
	for i := 0; i < handlerRef.Type().NumIn(); i++ {
		v, ok := inParams[handlerRef.Type().In(i)]
		if ok {
			params = append(params, v)
		} else {
			msg := fmt.Sprintf("cannot find inParam of [%v]", handlerRef.Type().In(i))
			return nil, errors.New(msg)
		}

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
func newType(fullRequestPath string, typ reflect.Type, c Context) (reflect.Value, error) {
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
			if lowCamelStr(pname) == lowCamelStr(field.Name) {
				isPathParam = true
				break
			}
		}
		if field.Name != "Body" {
			var value string
			if isPathParam {
				value = c.Param(field.Name)
			} else {
				queries := c.Request().URL().QueryParams()
				for k := range queries {
					if strings.ToLower(k) == strings.ToLower(field.Name) {
						value = c.Request().URL().QueryParam(k)
					}
				}
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

			if err := c.Bind(body); err != nil {
				return requestObj, err
			}

			if bodyType.Kind() == reflect.Ptr {
				requestObj.Elem().FieldByName("Body").Set(reflect.ValueOf(body))
			} else {
				requestObj.Elem().FieldByName("Body").Set(reflect.ValueOf(body).Elem())
			}
		}

	}
	err := validate.Struct(requestObj.Interface())
	return requestObj, err
}

func setValue(field reflect.Value, name string, value string) error {
	v, err := toTargeType(field.Type(), value)
	fmt.Printf("setValue [%s] -> %s(%v)\n", name, v, v.Type())
	field.Set(v)

	return err
}
