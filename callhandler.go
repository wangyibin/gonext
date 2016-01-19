package apidoc

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v8"
)

var validate *validator.Validate

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

// BuildEchoHandler func
func BuildEchoHandler(fullRequestPath string, handlers []interface{}) echo.HandlerFunc {
	inTypes, _, _ := validateChain(handlers)

	return func(c *echo.Context) error {
		// var requestObj reflect.Value
		var err error

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

		for _, h := range handlers {
			lastHandler = h
			out, err = callHanlder(h, inParams)
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
func newType(fullRequestPath string, typ reflect.Type, c *echo.Context) (reflect.Value, error) {
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
				queries := c.Request().URL.Query()
				for k := range c.Request().URL.Query() {
					if strings.ToLower(k) == strings.ToLower(field.Name) {
						value = queries.Get(k)
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

			c.Bind(body)
			fmt.Printf("%s\n", body)
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
