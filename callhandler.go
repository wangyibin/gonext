package gonext

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v8"
	"errors"
	"github.com/gorilla/schema"
	"time"
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
		StartAt := time.Now()

		var logError = func (err error) error {
			fmt.Printf("%4s | %3d [%.3fs] | %s\n", echoContext.Request().Method(),
				echoContext.Response().Status(), time.Now().Sub(StartAt).Seconds(),
				fullRequestPath)
			return err
		}
		var err error
		var c = NewGonextContextFromEcho(echoContext)
		inParams := make(map[reflect.Type]reflect.Value)
		inParams[reflect.TypeOf((*Context)(nil)).Elem()] = reflect.ValueOf(c)
		for _, inType := range inTypes {
			requestObj, err := newType(fullRequestPath, inType, c)
			if err != nil {
				return logError(err)
			}
			inParams[inType] = requestObj
		}

		var lastHandler interface{}
		var out []reflect.Value

		//fmt.Printf("call %s\n", fullRequestPath)
		//for inParamKey := range inParams {
		//	fmt.Printf("    in[%s]\n", inParamKey)
		//}
		for _, h := range handlers {
			lastHandler = h
			out, err = callHandler(h, inParams)
			if err != nil {
				return logError(err)
			}
		}
		if len(out) > 1 {
			return logError(fmt.Errorf("return more then one data value is not supported: %s", runtime.FuncForPC(reflect.ValueOf(lastHandler).Pointer()).Name()))
		} else if len(out) == 0 {
			return logError(c.NoContent(http.StatusOK))
		} else {
			return logError(c.JSON(http.StatusOK, out[0].Interface()))
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
			fmt.Println(msg)
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

	pathAndQueryParams := c.QueryParams()

	for _, name := range c.ParamNames() {
		pathAndQueryParams[name] = []string{c.Param(name)}
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(requestObj.Interface(), pathAndQueryParams)
	if err != nil {
		return requestObj, err
	}
	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)

		if field.Name == "Body" {
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
	err = validate.Struct(requestObj.Interface())
	return requestObj, err
}

//func setValue(field reflect.Value, name string, value string) error {
//	v, err := typeconv.ToTargetType(field.Type(), value)
//	fmt.Printf("setValue [%s] -> %s(%v)\n", name, v, v.Type())
//	field.Set(v)
//
//	return err
//}
