package apidoc

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Handler type
// type Handler interface{}

func validateChain(handlers []interface{}) (reflect.Type, reflect.Type, error) {
	var totalIns []reflect.Type
	var totalOuts []reflect.Type

	for _, h := range handlers {
		addTypes(h, &totalIns, &totalOuts)
	}

	uniqueIns := findUniqueTypes(totalIns, totalOuts)
	uniqueOuts := findUniqueTypes(totalOuts, totalIns)
	if len(uniqueIns) > 1 {
		return nil, nil, fmt.Errorf("more then one unique input type: %s", uniqueIns)
	}
	if len(uniqueOuts) > 1 {
		return nil, nil, fmt.Errorf("more then one unique output type: %s", uniqueOuts)
	}
	var uniqueIn reflect.Type
	var uniqueOut reflect.Type
	if len(uniqueIns) == 1 {
		uniqueIn = uniqueIns[0]
	}
	if len(uniqueOuts) == 1 {
		uniqueOut = uniqueOuts[0]
	}
	return uniqueIn, uniqueOut, nil
}

func addTypes(handler interface{}, totalIns *[]reflect.Type, totalOuts *[]reflect.Type) {
	if handler == nil {
		return
	}
	handlerType := reflect.TypeOf(handler)
	// fmt.Printf("handler type >> %s\n", handlerType)
	for i := 0; i < handlerType.NumIn(); i++ {
		if handlerType.In(i).String() != "*echo.Context" {
			*totalIns = append(*totalIns, handlerType.In(i))
		}
	}
	for i := 0; i < handlerType.NumOut(); i++ {
		if handlerType.Out(i).String() != "error" {
			*totalOuts = append(*totalOuts, handlerType.Out(i))
		}
	}
}
func getOperationID(inType reflect.Type, handlers []interface{}) string {
	// var operationHandler = handlers[len(handlers)-1]

	var fullName string
	for i := len(handlers) - 1; i >= 0; i-- {
		if isInTypeDefined(inType, handlers[i]) {
			fullName = runtime.FuncForPC(reflect.ValueOf(handlers[i]).Pointer()).Name()
		}
	}
	// if isInTypeDefined(inType, h3) {
	// 	operationHandler = h3
	// } else if isInTypeDefined(inType, h2) {
	// 	operationHandler = h2
	// } else if isInTypeDefined(inType, h1) {
	// 	operationHandler = h1
	// } else {
	// 	panic("getOperationID error happend")
	// }
	// fullName := runtime.FuncForPC(reflect.ValueOf(operationHandler).Pointer()).Name()
	arr := strings.Split(fullName, ".")
	return arr[len(arr)-1]
}

func isInTypeDefined(inType reflect.Type, handler interface{}) bool {
	if handler == nil {
		return false
	}
	if inType == nil {
		return true
	}
	handlerType := reflect.TypeOf(handler)
	for i := 0; i < handlerType.NumIn(); i++ {
		if handlerType.In(i).String() == inType.String() {
			return true
		}
	}
	return false
}

func findUniqueTypes(source []reflect.Type, compareWith []reflect.Type) []reflect.Type {
	var uniques []reflect.Type
	for _, in := range source {
		isUnique := true
		for _, out := range compareWith {
			if in.String() == out.String() {
				isUnique = false
			}
		}
		if isUnique {
			uniques = append(uniques, in)
		}
	}
	return uniques
}
