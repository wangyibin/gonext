package apidoc

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo"
)

// Router struct
type Router struct {
	e         *echo.Echo
	resources []*Resource
}

// Resource struct
type Resource struct {
	method      string
	path        string
	description string
	pathParams  []Param
	queryParams []Param
	requestBody reflect.Type
	operationID string
}

// Param struct
type Param struct {
	Name     string
	Type     string
	Required bool
}

func (p *Param) String() string {
	return fmt.Sprintf("name[%s] type:%s, required:%t\n", p.Name, p.Type, p.Required)
}

// NewRouter func
func NewRouter(e *echo.Echo) *Router {
	r := &Router{e: e}
	e.Get("/api-docs", getdoc(r))
	return r
}

// Get func
func (r *Router) Get(path string, h interface{}) {
	r.mountToAPIDocs("GET", path, h).mountToRouter(r.e, h)
}

// Post func
func (r *Router) Post(path string, h interface{}) {
	r.mountToAPIDocs("POST", path, h).mountToRouter(r.e, h)
}

func (r *Router) mountToAPIDocs(method string, path string, h interface{}) *Resource {
	resource := &Resource{method: method, path: path, operationID: runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()}

	handlerType := reflect.TypeOf(h)
	requestType := handlerType.In(0)
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}

	pnames := pathNames(path)
	fmt.Printf("%s %s\n\trequest type %v\n", method, path, requestType)
	for i := 0; i < requestType.NumField(); i++ {
		typeField := requestType.Field(i)

		if strings.ToUpper(typeField.Name) != "BODY" {
			param := Param{Name: typeField.Name, Type: typeField.Type.String(), Required: typeField.Type.Kind() != reflect.Ptr}
			if containsIgnoreCase(pnames, typeField.Name) {
				resource.pathParams = append(resource.queryParams, param)
				fmt.Printf("\tPath Params %s", param.String())
			} else {
				resource.queryParams = append(resource.queryParams, param)
				fmt.Printf("\tQuery Params %s", param.String())
			}
		} else {
			if typeField.Type.Kind() == reflect.Ptr {
				resource.requestBody = typeField.Type.Elem()
			} else {
				resource.requestBody = typeField.Type
			}
			fmt.Printf("\tRequestBody[%s]\n", resource.requestBody)
		}
	}

	r.resources = append(r.resources, resource)
	return resource
}

func containsIgnoreCase(s []string, e string) bool {
	for _, a := range s {
		if strings.ToUpper(a) == strings.ToUpper(e) {
			return true
		}
	}
	return false
}
func (r *Resource) mountToRouter(e *echo.Echo, h interface{}) {
	e.Match([]string{r.method}, r.path, func(c *echo.Context) error {
		handlerType := reflect.TypeOf(h)
		requestType := handlerType.In(0)
		if requestType.Kind() == reflect.Ptr {
			requestType = requestType.Elem()
		}
		requestObj := reflect.New(requestType)

		for _, queryParam := range r.queryParams {
			paramName := queryParam.Name
			value := c.Query(lowCamelStr(paramName))

			setValue(requestObj.Elem().FieldByName(paramName), paramName, value)
		}
		for _, pathParam := range r.pathParams {
			paramName := pathParam.Name
			value := c.Param(lowCamelStr(pathParam.Name))
			setValue(requestObj.Elem().FieldByName(paramName), paramName, value)
		}
		if r.requestBody != nil {
			body := reflect.New(r.requestBody).Interface()
			c.Bind(body)
			fmt.Printf("%s\n", body)
			requestObj.Elem().FieldByName("Body").Set(reflect.ValueOf(body))
		}
		result := reflect.ValueOf(h).Call([]reflect.Value{requestObj, reflect.ValueOf(c)})[0].Interface()
		if result != nil {
			return result.(error)
		}
		return nil
	})
}
func setValue(field reflect.Value, name string, value string) error {
	v, err := toTargeType(field.Type(), value)
	fmt.Printf("setValue [%s] -> %s(%v)\n", name, v, v.Type())
	field.Set(v)

	return err
}

func pathNames(path string) []string {
	pnames := []string{} // Param names
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)
		} else if path[i] == '*' {
			pnames = append(pnames, "_*")
		}
	}
	return pnames
}

func lowCamelStr(str string) string {
	return strings.ToLower(string(str[0])) + string(str[1:])
}

// func entityName(v interface{}) string {
// 	fullName := reflect.TypeOf(v)
// 	arr := strings.Split(fullName.String(), ".")
// 	obj := arr[len(arr)-1]
// 	return strings.ToUpper(string(obj[0])) + string(obj[1:])
// }
func (r *Resource) toJSON() map[string]interface{} {
	var parameters []map[string]interface{}
	if r.requestBody != nil {
		parameters = append(parameters, map[string]interface{}{
			"in":       "body",
			"name":     "body",
			"required": true,
			"schema": map[string]string{
				"$ref": "#/definitions/" + r.requestBody.Name(),
			},
		})
	}
	return map[string]interface{}{
		strings.ToLower(r.method): map[string]interface{}{
			"tags":        []string{"pet"},
			"summary":     "Add a new pet to the store",
			"description": r.description,
			"produces":    []string{"application/json"},
			"consumes":    []string{"application/json"},
			"operationId": r.operationID,
			"parameters":  parameters,
			"response": map[string]interface{}{
				"405": map[string]interface{}{
					"description": "Invalid input",
				},
			},
		},
	}
}
